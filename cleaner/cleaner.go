package cleaner

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

type Cleaner struct {
	logger *logrus.Entry
	ticker *time.Ticker
	hlsDir string
}

func NewCleaner(ctx context.Context, hlsDir string) (*Cleaner, error) {
	return &Cleaner{
		logger: ctxlogrus.Extract(ctx).WithField("system", "cleaner"),
		ticker: time.NewTicker(time.Second * 60),
		hlsDir: hlsDir,
	}, nil
}

func (c *Cleaner) Start() {
	for range c.ticker.C {
		err := c.cleanup()
		if err != nil {
			c.logger.WithError(err).Error("failed to cleanup")
		}
	}
}

func (c *Cleaner) Stop() {
	c.ticker.Stop()
}

func (c *Cleaner) cleanup() error {
	if _, err := os.Stat(c.hlsDir); os.IsNotExist(err) {
		return nil
	}

	files, err := ioutil.ReadDir(c.hlsDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if f.ModTime().Before(time.Now().Add(time.Hour * -12)) {
				path := filepath.Join(c.hlsDir, f.Name())
				c.logger.WithField("path", path).Info("removing")
				err := os.RemoveAll(path)
				if err != nil {
					c.logger.WithError(err).Error("failed to remove")
					continue
				}
			}
		}
	}

	return nil
}
