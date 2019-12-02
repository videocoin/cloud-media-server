package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string        `envconfig:"-"`
	Version string        `envconfig:"-"`
	Logger  *logrus.Entry `envconfig:"-"`

	RPCAddr         string `default:"0.0.0.0:5090" envconfig:"RPC_ADDR"`
	AuthTokenSecret string `default:"" envconfig:"AUTH_TOKEN_SECRET"`
	UsersRPCAddr    string `default:"0.0.0.0:5000" envconfig:"USERS_RPC_ADDR"`
	StreamsRPCAddr  string `default:"0.0.0.0:5102" envconfig:"STREAMS_RPC_ADDR"`

	MediaServerHost string `default:"127.0.0.1" envconfig:"MEDIASERVER_HOST"`
	MediaServerPort int    `default:"6000" envconfig:"MEDIASERVER_PORT"`

	RTMPURL string `required:"true" envconfig:"RTMP_URL"`
}
