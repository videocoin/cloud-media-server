module github.com/videocoin/mediaserver

go 1.12

require (
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/notedit/gstreamer-rtmp v0.0.0-20181226050148-9295bf2f2ca8
	github.com/notedit/media-server-go v0.1.14
	github.com/notedit/sdp v0.0.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/videocoin/cloud-api v0.2.15
	github.com/videocoin/cloud-pkg v0.0.6
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135 // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/go-playground/validator.v9 v9.30.2
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
)

replace github.com/videocoin/cloud-api => ../cloud-api

replace github.com/videocoin/cloud-pkg => ../cloud-pkg
