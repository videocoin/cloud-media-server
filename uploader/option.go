package uploader

import (
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-media-server/datastore"
)

type Option func(*Uploader)

func WithLogger(logger *logrus.Entry) Option {
	return func(d *Uploader) {
		d.logger = logger
	}
}

func WithDatastore(ds datastore.Datastore) Option {
	return func(d *Uploader) {
		d.ds = ds
	}
}

func WithServiceClient(sc *clientv1.ServiceClient) Option {
	return func(d *Uploader) {
		d.sc = sc
	}
}

func WithDataDir(dataDir string) Option {
	return func(d *Uploader) {
		d.dataDir = dataDir
	}
}

func WithGDriveKey(key string) Option {
	return func(d *Uploader) {
		d.gdriveKey = key
	}
}
