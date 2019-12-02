package rpc

import (
	"context"
	"net"

	"github.com/videocoin/mediaserver/mediacore"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/mediaserver/v1"
	"github.com/videocoin/cloud-api/rpc"
	streamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/auth"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type RpcServerOpts struct {
	Logger          *logrus.Entry
	Addr            string
	AuthTokenSecret string
	Users           usersv1.UserServiceClient
	Streams         streamsv1.StreamsServiceClient
	MS              *mediacore.MediaServer
}

type RpcServer struct {
	logger          *logrus.Entry
	addr            string
	authTokenSecret string
	rtmpURL         string
	grpc            *grpc.Server
	listen          net.Listener
	users           usersv1.UserServiceClient
	streams         streamsv1.StreamsServiceClient
	validator       *requestValidator
	ms              *mediacore.MediaServer
}

func NewRpcServer(opts *RpcServerOpts) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		logger:          opts.Logger.WithField("system", "rpc"),
		validator:       newRequestValidator(),
		addr:            opts.Addr,
		authTokenSecret: opts.AuthTokenSecret,
		grpc:            grpcServer,
		listen:          listen,
		users:           opts.Users,
		streams:         opts.Streams,
		ms:              opts.MS,
	}

	v1.RegisterMediaServerServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RpcServer) authenticate(ctx context.Context) (string, context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "authenticate")
	defer span.Finish()

	ctx = auth.NewContextWithSecretKey(ctx, s.authTokenSecret)
	ctx, jwtToken, err := auth.AuthFromContext(ctx)
	if err != nil {
		s.logger.Warningf("failed to auth from context: %s", err)
		return "", ctx, rpc.ErrRpcUnauthenticated
	}

	tokenType, ok := auth.TypeFromContext(ctx)
	if ok {
		if usersv1.TokenType(tokenType) == usersv1.TokenTypeAPI {
			_, err := s.users.GetApiToken(context.Background(), &usersv1.ApiTokenRequest{Token: jwtToken})
			if err != nil {
				s.logger.Errorf("failed to get api token: %s", err)
				return "", ctx, rpc.ErrRpcUnauthenticated
			}
		}
	}

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		s.logger.Warningf("failed to get user id from context: %s", err)
		return "", ctx, rpc.ErrRpcUnauthenticated
	}

	return userID, ctx, nil
}
