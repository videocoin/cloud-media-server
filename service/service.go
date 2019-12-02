package service

import (
	streamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/mediaserver/mediacore"
	"github.com/videocoin/mediaserver/rpc"
)

type Service struct {
	cfg *Config
	rpc *rpc.RpcServer
	ms  *mediacore.MediaServer
}

func NewService(cfg *Config) (*Service, error) {
	conn, err := grpcutil.Connect(cfg.UsersRPCAddr, cfg.Logger.WithField("system", "userscli"))
	if err != nil {
		return nil, err
	}
	users := usersv1.NewUserServiceClient(conn)

	conn, err = grpcutil.Connect(cfg.StreamsRPCAddr, cfg.Logger.WithField("system", "pstreamscli"))
	if err != nil {
		return nil, err
	}
	streams := streamsv1.NewStreamsServiceClient(conn)

	mediaServer, err := mediacore.NewMediaServer(
		cfg.MediaServerHost,
		cfg.MediaServerPort,
		cfg.RTMPURL,
		cfg.Logger.WithField("system", "mediaserver"),
	)
	if err != nil {
		return nil, err
	}

	rpcConfig := &rpc.RpcServerOpts{
		Logger:  cfg.Logger,
		Addr:    cfg.RPCAddr,
		Streams: streams,
		Users:   users,
		MS:      mediaServer,
	}

	rpcServer, err := rpc.NewRpcServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg: cfg,
		rpc: rpcServer,
		ms:  mediaServer,
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.rpc.Start()
	go s.ms.Start()
	return nil
}

func (s *Service) Stop() error {
	s.ms.Stop()
	return nil
}
