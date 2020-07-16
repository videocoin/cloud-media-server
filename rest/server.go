package rest

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	otecho "github.com/opentracing-contrib/echo"
	"github.com/opentracing/opentracing-go"
	echologrus "github.com/plutov/echo-logrus"
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-media-server/datastore"
	"github.com/videocoin/cloud-media-server/downloader"
	"github.com/videocoin/cloud-media-server/eventbus"
	"github.com/videocoin/cloud-media-server/mediacore/splitter"
	"github.com/videocoin/cloud-media-server/nginxrtmp"
	"net/http"
)

type Server struct {
	addr            string
	authTokenSecret string
	bucket          string
	logger          *logrus.Entry
	e               *echo.Echo
	downloader      *downloader.Downloader
	splitter        *splitter.Splitter
	nginxRTMPHook   *nginxrtmp.Hook
	sc              *clientv1.ServiceClient
	ds              datastore.Datastore
	eb              *eventbus.EventBus
	gs              *storage.Client
	bh              *storage.BucketHandle
}

func NewServer(ctx context.Context, opts ...Option) (*Server, error) {
	s := &Server{
		logger: ctxlogrus.Extract(ctx).WithField("system", "api"),
		e:      echo.New(),
	}

	echologrus.Logger = s.logger.Logger

	s.e.HideBanner = true
	s.e.HidePort = true
	s.e.DisableHTTP2 = true
	s.e.Logger = echologrus.GetEchoLogger()

	for _, o := range opts {
		err := o(s)
		if err != nil {
			return nil, err
		}
	}

	err := s.route(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) route(ctx context.Context) error {
	s.e.Use(otecho.Middleware("mediaserver"))
	s.e.Use(middleware.CORS())
	s.e.Use(echologrus.Hook())

	s.e.GET("/healthz", s.getHealth)

	v1 := s.e.Group("/api/v1/")

	// Uploader
	uploaderV1 := v1.Group("upload/")
	uploaderV1.Use(middleware.JWT([]byte(s.authTokenSecret)))
	uploaderV1.POST("local/:id", s.uploadFromFile)
	uploaderV1.POST("url/:id", s.uploadFromURL)
	uploaderV1.GET("url/:id", s.getUploadInfo)

	// Syncer
	v1.POST("sync", s.sync)

	// Nginx RTMP hooks
	hookConfig := &nginxrtmp.HookConfig{Prefix: "/nginx-rtmp/hooks"}
	hookLogger := ctxlogrus.ToContext(ctx, s.logger.WithField("system", "nginx-rtmp"))
	nginxRTMPHook, err := nginxrtmp.NewHook(hookLogger, s.e, hookConfig, s.sc)
	if err != nil {
		return err
	}
	s.nginxRTMPHook = nginxRTMPHook

	return nil
}

func (s *Server) Start() error {
	return s.e.Start(s.addr)
}

func (s *Server) Stop() error {
	return s.e.Shutdown(context.Background())
}

func (s *Server) E() *echo.Echo {
	return s.e
}

func (s *Server) getHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{"alive": true})
}

func (s *Server) authenticate(ctx echo.Context) (string, error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	span, _ := opentracing.StartSpanFromContext(ctx.Request().Context(), "authenticate")
	defer span.Finish()

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}

	return userID, nil
}
