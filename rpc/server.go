package rpc

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	v1 "github.com/videocoin/cloud-api/mediaserver/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/mediaserver/mediacore/webrtc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

type ServerOpts struct {
	Addr            string
	AuthTokenSecret string
	Bucket          string
	SC              *clientv1.ServiceClient
	WebRTCStreamer  *webrtc.Streamer
}

type Server struct {
	logger          *logrus.Entry
	addr            string
	authTokenSecret string
	bucket          string
	grpc            *grpc.Server
	listen          net.Listener
	validator       *requestValidator
	sc              *clientv1.ServiceClient
	webRtcStreamer  *webrtc.Streamer
	bh              *storage.BucketHandle
}

func NewServer(ctx context.Context, opts *ServerOpts) (*Server, error) {
	grpcOpts := grpcutil.DefaultServerOpts(ctxlogrus.Extract(ctx))
	grpcServer := grpc.NewServer(grpcOpts...)

	healthService := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthService)

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	gscli, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	bh := gscli.Bucket(opts.Bucket)
	_, err = bh.Attrs(context.Background())
	if err != nil {
		return nil, err
	}

	rpcServer := &Server{
		logger:          ctxlogrus.Extract(ctx).WithField("system", "rpc"),
		validator:       newRequestValidator(),
		addr:            opts.Addr,
		authTokenSecret: opts.AuthTokenSecret,
		grpc:            grpcServer,
		listen:          listen,
		sc:              opts.SC,
		webRtcStreamer:  opts.WebRTCStreamer,
		bh:              bh,
	}

	v1.RegisterMediaServerServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *Server) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *Server) Stop() error {
	s.grpc.GracefulStop()
	return nil
}
