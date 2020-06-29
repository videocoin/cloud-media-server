package service

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/sirupsen/logrus"
	streamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-media-server/mediacore"
	"github.com/videocoin/cloud-media-server/rpc"
	"github.com/videocoin/cloud-pkg/grpcutil"
)

type Service struct {
	cfg    *Config
	logger *logrus.Entry
	rpc    *rpc.Server
	ms     *mediacore.MediaServer
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

	gscli, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	bh := gscli.Bucket(cfg.Bucket)
	_, err = bh.Attrs(context.Background())
	if err != nil {
		return nil, err
	}

	rpcConfig := &rpc.ServerOpts{
		Logger:          cfg.Logger,
		Addr:            cfg.RPCAddr,
		AuthTokenSecret: cfg.AuthTokenSecret,
		Streams:         streams,
		Users:           users,
		MS:              mediaServer,
		GS:              gscli,
		BH:              bh,
	}

	rpcServer, err := rpc.NewServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:    cfg,
		logger: cfg.Logger,
		rpc:    rpcServer,
		ms:     mediaServer,
	}

	return svc, nil
}

func (s *Service) Start() {
	go func() {
		err := s.rpc.Start()
		if err != nil {
			s.logger.WithError(err).Fatal("failed to start rpc server")
		}
	}()

	go func() {
		err := s.ms.Start()
		if err != nil {
			s.logger.WithError(err).Fatal("failed to start media server")
		}
	}()
}

func (s *Service) Stop() error {
	return s.ms.Stop()
}
