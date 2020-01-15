FROM gcr.io/videocoin-network/mediaserver-go:v0.1.0-develop-f87201a as builder

RUN apt-get update
RUN apt-get install -y opus-tools libogg0 libopus0 libopus-dev libopusfile-dev

RUN go build 

WORKDIR /go/src/github.com/videocoin/mediaserver

ADD ./ ./

RUN make build

FROM ubuntu
RUN apt-get update
RUN apt-get install -y ca-certificates \
    libgstreamer-plugins-base1.0-dev \
    libgstreamer1.0-0 \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    gstreamer1.0-tools \
    gstreamer1.0-x \
    gstreamer1.0-pulseaudio \
    opus-tools \
    libogg0 \
    libopus0 \
    libopus-dev \
    libopusfile-dev

COPY --from=builder /go/src/github.com/videocoin/mediaserver/bin/mediaserver /mediaserver

WORKDIR /

CMD ["/mediaserver"]
