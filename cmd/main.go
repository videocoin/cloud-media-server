package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/kelseyhightower/envconfig"
	pkglogger "github.com/videocoin/cloud-pkg/logger"
	"github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/cloud-media-server/app"
)

var (
	Name    string = "mediaserver"
	Version string = "dev"
)

func main() {
	logger := pkglogger.NewLogrusLogger(Name, Version, nil)

	closer, err := tracer.NewTracer(Name)
	if err != nil {
		logger.WithError(err).Info("failed to new tracer")
	} else {
		defer closer.Close()
	}

	cfg := &app.Config{
		Name:    Name,
		Version: Version,
	}

	err = envconfig.Process(Name, cfg)
	if err != nil {
		logger.WithError(err).Fatal("failed to process config")
	}

	ctx := ctxlogrus.ToContext(context.Background(), logger)
	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.WithError(err).Fatal("failed to create app")
	}

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	errCh := make(chan error, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		logger.WithField("signal", sig.String()).Infof("received signal")
		exit <- true
	}()

	logger.Info("starting")
	go app.Start(errCh)

	select {
	case <-exit:
		break
	case err := <-errCh:
		if err != nil {
			logger.WithError(err).Error("failed to start app")
		}
		break
	}

	logger.Info("stopping")
	err = app.Stop()
	if err != nil {
		logger.WithError(err).Error("failed to stop app")
		return
	}

	logger.Info("stopped")
}
