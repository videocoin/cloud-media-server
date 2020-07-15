package uploader

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/mediaserver/datastore"
	"github.com/videocoin/mediaserver/mediacore/splitter"
	"github.com/videocoin/mediaserver/downloader"
)

type Uploader struct {
	logger     *logrus.Entry
	stop       chan bool
	dataDir string
	gdriveKey string
	downloader *downloader.Downloader
	splitter   *splitter.Splitter
	sc *clientv1.ServiceClient
	ds datastore.Datastore
}

func NewUploader(ctx context.Context, opts ...Option) (*Uploader, error) {
	uploader := &Uploader{
		logger: ctxlogrus.Extract(ctx).WithField("system", "uploader"),
		stop:   make(chan bool, 1),
	}

	for _, o := range opts {
		o(uploader)
	}

	splitter, err := splitter.NewSplitter(ctx, splitter.WithOutputDir(uploader.dataDir))
	if err != nil {
		return nil, err
	}
	uploader.splitter = splitter

	downloaderCtx := downloader.NewContextWithGDriveKey(ctx, uploader.gdriveKey)
	downloader, err := downloader.NewDownloader(
		downloaderCtx,
		uploader.dataDir,
		downloader.WithDatastore(uploader.ds),
	)
	if err != nil {
		return nil, err
	}
	uploader.downloader = downloader

	return uploader, nil
}

func (u *Uploader) dispatch() {
	go func() {
		for outputFile := range u.downloader.OutputCh {
			if outputFile != nil && u.splitter != nil {
				logger := u.logger.WithField("stream_id", outputFile.StreamID)
				logger.Info("received output file from downloader")

				ctx := ctxlogrus.ToContext(context.Background(), logger)

				if outputFile.Error != nil {
					logger.WithError(outputFile.Error).Info("failed to download")
					if u.sc != nil && u.sc.Streams != nil {
						u.stopStream(ctx, outputFile.StreamID)
					}
					continue
				}

				go func() {
					logger.Info("sending file to splitter")
					u.splitter.InputCh <- &splitter.MediaFile{
						StreamID: outputFile.StreamID,
						Path:     outputFile.Path,
					}
				}()
			}
		}
	}()

	go func() {
		for mf := range u.splitter.OutputCh {
			logger := u.logger.WithField("stream_id", mf.StreamID).WithField("path", mf.Path)
			logger.Info("recieved output from splitter")

			ctx := ctxlogrus.ToContext(context.Background(), logger)

			if mf.Error != nil {
				logger.WithError(mf.Error).Info("failed to split")
				if u.sc != nil && u.sc.Streams != nil {
					u.stopStream(ctx, mf.StreamID)
				}
				continue
			}

			u.logger.
				WithField("stream_id", mf.StreamID).
				WithField("path", mf.Path).
				Info("file has been splitted")

			if u.sc != nil && u.sc.Streams != nil {
				logger.Info("stream publishing")

				streamReq := &pstreamsv1.StreamRequest{
					Id:       mf.StreamID,
					Duration: mf.Duration,
				}
				_, err := u.sc.Streams.Publish(context.Background(), streamReq)
				if err != nil {
					logger.WithError(err).Error("failed to publish stream")
					u.stopStream(ctx, mf.StreamID)
					continue
				}
			}
		}
	}()

	u.stop <- true
}

func (u *Uploader) stopStream(ctx context.Context, streamID string) {
	logger := ctxlogrus.Extract(ctx)
	streamReq := &pstreamsv1.StreamRequest{Id: streamID}
	_, err := u.sc.Streams.Stop(context.Background(), streamReq)
	if err != nil {
		logger.WithError(err).Info("failed to stop stream")
	}
}

func (u *Uploader) Start(errCh chan error) {
	go func() {
		u.logger.WithField("dst", u.dataDir).Info("starting splitter")
		err := u.splitter.Start()
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		u.logger.WithField("dst", u.dataDir).Info("starting downloader")
		err := u.downloader.Start()
		if err != nil {
			errCh <- err
		}
	}()

	u.dispatch()
}

func (u *Uploader) Stop() error {
	close(u.stop)

	u.logger.Info("stopping splitter")
	err := u.splitter.Stop()
	if err != nil {
		return err
	}

	u.logger.Info("stopping downloader")
	err = u.downloader.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (u *Uploader) Downloader() *downloader.Downloader {
	return u.downloader
}

func (u *Uploader) Splitter() *splitter.Splitter {
	return u.splitter
}