package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string        `envconfig:"-"`
	Version string        `envconfig:"-"`
	Logger  *logrus.Entry `envconfig:"-"`

	RPCAddr         string `envconfig:"RPC_ADDR" default:"0.0.0.0:5090"`
	AuthTokenSecret string `envconfig:"AUTH_TOKEN_SECRET" default:""`
	UsersRPCAddr    string `envconfig:"USERS_RPC_ADDR" default:"0.0.0.0:5000"`
	StreamsRPCAddr  string `envconfig:"STREAMS_RPC_ADDR" default:"0.0.0.0:5102"`

	MediaServerHost string `envconfig:"MEDIASERVER_HOST" default:"127.0.0.1"`
	MediaServerPort int    `envconfig:"MEDIASERVER_PORT" default:"6000"`

	RTMPURL string `envconfig:"RTMP_URL" required:"true"`
	Bucket  string `envconfig:"BUCKET" required:"true"`
}
