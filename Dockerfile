FROM registry.videocoin.net/cloud/mediaserver-base:v0.1

ADD ./deploy/nginx-rtmp /usr/local/nginx-rtmp

WORKDIR /root/go/src/github.com/videocoin/mediaserver
ADD ./ ./
RUN make build

RUN cp bin/mediaserver /mediaserver-linux-amd64

WORKDIR /

ENTRYPOINT ["/mediaserver-linux-amd64"]
