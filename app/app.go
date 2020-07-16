package app

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-media-server/cleaner"
	"github.com/videocoin/cloud-media-server/datastore"
	"github.com/videocoin/cloud-media-server/eventbus"
	"github.com/videocoin/cloud-media-server/mediacore/webrtc"
	"github.com/videocoin/cloud-media-server/rest"
	"github.com/videocoin/cloud-media-server/rpc"
	"github.com/videocoin/cloud-media-server/uploader"
)

type App struct {
	cfg      *Config
	logger   *logrus.Entry
	rest     *rest.Server
	rpc      *rpc.Server
	uploader *uploader.Uploader
	cleaner  *cleaner.Cleaner
	ds       datastore.Datastore
	eb       *eventbus.EventBus
	sc       *clientv1.ServiceClient
	stop     chan bool
}

func NewApp(ctx context.Context, cfg *Config) (*App, error) {
	sc, err := clientv1.NewServiceClientFromEnvconfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	ds, err := datastore.NewDatastore(ctx, cfg.RedisURI)
	if err != nil {
		return nil, err
	}

	eb, err := eventbus.NewEventBus(
		ctx,
		cfg.MQURI,
		eventbus.WithName(cfg.Name),
		eventbus.WithBucket(cfg.Bucket),
	)
	if err != nil {
		return nil, err
	}

	uploader, err := uploader.NewUploader(
		ctx,
		uploader.WithDataDir(cfg.DataDir),
		uploader.WithGDriveKey(cfg.GDriveKey),
		uploader.WithDatastore(ds),
		uploader.WithServiceClient(sc),
	)
	if err != nil {
		return nil, err
	}

	restServerOpts := []rest.Option{
		rest.WithAddr(cfg.Addr),
		rest.WithAuthTokenSecret(cfg.AuthTokenSecret),
		rest.WithDownloader(uploader.Downloader()),
		rest.WithServiceClient(sc),
		rest.WithDatastore(ds),
		rest.WithSplitter(uploader.Splitter()),
		rest.WithBucket(cfg.Bucket),
		rest.WithEventBus(eb),
	}
	rest, err := rest.NewServer(ctx, restServerOpts...)
	if err != nil {
		return nil, err
	}

	webRtcStreamer, err := webrtc.NewStreamer(ctx, cfg.WebRTCServerHost, cfg.WebRTCServerPort, cfg.RTMPURL)
	if err != nil {
		return nil, err
	}

	rpcServerOpts := &rpc.ServerOpts{
		Addr:            cfg.RPCAddr,
		AuthTokenSecret: cfg.AuthTokenSecret,
		Bucket: cfg.Bucket,
		SC:              sc,
		WebRTCStreamer:  webRtcStreamer,
	}
	rpc, err := rpc.NewServer(ctx, rpcServerOpts)
	if err != nil {
		return nil, err
	}

	cleaner, err := cleaner.NewCleaner(ctx, cfg.DataDir)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:      cfg,
		logger:   ctxlogrus.Extract(ctx),
		stop:     make(chan bool, 1),
		rest:     rest,
		rpc:      rpc,
		ds:       ds,
		eb:       eb,
		sc:       sc,
		uploader: uploader,
		cleaner:  cleaner,
	}, nil
}

func (s *App) Start(errCh chan error) {
	go func() {
		s.logger.WithField("addr", s.cfg.Addr).Info("starting http server")
		errCh <- s.rest.Start()
	}()

	go func() {
		s.logger.WithField("addr", s.cfg.RPCAddr).Info("starting rpc server")
		errCh <- s.rpc.Start()
	}()

	go func() {
		s.logger.WithField("data_dir", s.cfg.DataDir).Info("starting uploader")
		s.uploader.Start(errCh)
	}()

	go func() {
		s.logger.WithField("data_dir", s.cfg.DataDir).Info("starting cleaner")
		s.cleaner.Start()
	}()

	go func() {
		s.logger.WithField("dburi", s.cfg.RedisURI).Info("starting datastore")
		err := s.ds.Start()
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		s.logger.WithField("mquri", s.cfg.MQURI).Info("starting eventbus")
		err := s.eb.Start()
		if err != nil {
			errCh <- err
		}
	}()
}

func (s *App) Stop() error {
	s.logger.Info("stopping api server")
	err := s.rest.Stop()
	if err != nil {
		return err
	}

	s.logger.Info("stopping rpc server")
	err = s.rpc.Stop()
	if err != nil {
		return err
	}

	s.logger.Info("stopping uploader")
	err = s.uploader.Stop()
	if err != nil {
		return err
	}

	s.logger.Info("stopping cleaner")
	s.cleaner.Stop()

	s.logger.Info("stopping datastore")
	err = s.ds.Stop()
	if err != nil {
		return err
	}

	s.logger.Info("stopping eventbus")
	err = s.eb.Stop()
	if err != nil {
		return err
	}

	s.stop <- true

	return nil
}
