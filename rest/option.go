package rest

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/mediaserver/datastore"
	"github.com/videocoin/mediaserver/downloader"
	"github.com/videocoin/mediaserver/eventbus"
	"github.com/videocoin/mediaserver/mediacore/splitter"
)

type Option func(*Server) error

func WithAddr(addr string) Option {
	return func(s *Server) error {
		s.addr = addr
		return nil
	}
}

func WithLogger(logger *logrus.Entry) Option {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

func WithDownloader(d *downloader.Downloader) Option {
	return func(s *Server) error {
		s.downloader = d
		return nil
	}
}

func WithAuthTokenSecret(secret string) Option {
	return func(s *Server) error {
		s.authTokenSecret = secret
		return nil
	}
}

func WithServiceClient(sc *clientv1.ServiceClient) Option {
	return func(s *Server) error {
		s.sc = sc
		return nil
	}
}

func WithDatastore(ds datastore.Datastore) Option {
	return func(s *Server) error {
		s.ds = ds
		return nil
	}
}

func WithSplitter(splitter *splitter.Splitter) Option {
	return func(s *Server) error {
		s.splitter = splitter
		return nil
	}
}

func WithBucket(name string) Option {
	return func(s *Server) error {
		s.bucket = name

		gs, err := storage.NewClient(context.Background())
		if err != nil {
			return err
		}
		s.gs = gs

		bh := gs.Bucket(s.bucket)
		_, err = bh.Attrs(context.Background())
		if err != nil {
			return err
		}
		s.bh = bh

		return nil
	}
}

func WithEventBus(eb *eventbus.EventBus) Option {
	return func(s *Server) error {
		s.eb = eb
		return nil
	}
}
