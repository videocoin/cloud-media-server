package rpc

import (
	"cloud.google.com/go/storage"
	"net"

	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/mediaserver/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-media-server/mediacore"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerOpts struct {
	Logger          *logrus.Entry
	Addr            string
	AuthTokenSecret string
	Users           usersv1.UserServiceClient
	Streams         streamsv1.StreamsServiceClient
	MS              *mediacore.MediaServer
	GS              *storage.Client
	BH              *storage.BucketHandle
}

type Server struct {
	logger          *logrus.Entry
	addr            string
	authTokenSecret string
	grpc            *grpc.Server
	listen          net.Listener
	users           usersv1.UserServiceClient
	streams         streamsv1.StreamsServiceClient
	validator       *requestValidator
	ms              *mediacore.MediaServer
	gs              *storage.Client
	bh              *storage.BucketHandle
}

func NewServer(opts *ServerOpts) (*Server, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)

	healthService := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthService)

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	rpcServer := &Server{
		logger:          opts.Logger.WithField("system", "rpc"),
		validator:       newRequestValidator(),
		addr:            opts.Addr,
		authTokenSecret: opts.AuthTokenSecret,
		grpc:            grpcServer,
		listen:          listen,
		users:           opts.Users,
		streams:         opts.Streams,
		ms:              opts.MS,
		gs:              opts.GS,
		bh:              opts.BH,
	}

	v1.RegisterMediaServerServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *Server) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}
