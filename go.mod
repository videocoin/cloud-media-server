module github.com/videocoin/cloud-media-server

go 1.14

require (
	cloud.google.com/go/storage v1.6.0
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/googleapis/gax-go v1.0.3 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/notedit/media-server-go v0.1.14
	github.com/notedit/sdp v0.0.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/videocoin/cloud-api v0.2.15
	github.com/videocoin/cloud-pkg v0.0.6
	go.opencensus.io v0.22.4 // indirect
	google.golang.org/api v0.28.0 // indirect
	google.golang.org/grpc v1.28.0
	gopkg.in/go-playground/validator.v9 v9.30.2
	gopkg.in/hraban/opus.v2 v2.0.0-20191117073431-57179dff69a6 // indirect
)

replace github.com/videocoin/cloud-api => ../cloud-api

replace github.com/videocoin/cloud-pkg => ../cloud-pkg
