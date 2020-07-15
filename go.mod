module github.com/videocoin/mediaserver

go 1.14

require (
	cloud.google.com/go/storage v1.10.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/grafov/m3u8 v0.11.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/notedit/media-server-go v0.1.14
	github.com/notedit/sdp v0.0.1
	github.com/opentracing-contrib/echo v0.0.0-20190807091611-5fe2e1308f06
	github.com/opentracing/opentracing-go v1.1.0
	github.com/plutov/echo-logrus v1.0.0
	github.com/sirupsen/logrus v1.6.0
	github.com/streadway/amqp v1.0.0
	github.com/vansante/go-ffprobe v1.1.0
	github.com/videocoin/cloud-api v1.0.0
	github.com/videocoin/cloud-pkg v1.0.0
	google.golang.org/api v0.28.0
	google.golang.org/grpc v1.30.0
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/redis.v5 v5.2.9
)

replace (
	github.com/videocoin/cloud-api => "../cloud-api"
	github.com/videocoin/cloud-pkg => "../cloud-pkg"
)
