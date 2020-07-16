FROM registry.videocoin.net/cloud/mediaserver-base:1.0

ADD ./deploy/nginx-rtmp /usr/local/nginx-rtmp

WORKDIR /root/go/src/github.com/videocoin/cloud-media-server
ADD ./ ./
RUN make build

RUN cp bin/mediaserver /mediaserver-linux-amd64

WORKDIR /

ENTRYPOINT ["/mediaserver-linux-amd64"]
