FROM registry.videocoin.net/cloud/mediaserver-base:1.0

RUN apt update
RUN apt install -y libglu1-mesa libswscale4 libavformat57 libavdevice57 libavcodec57 libjpeg62 libmad0
RUN wget http://download.tsi.telecom-paristech.fr/gpac/release/1.0.1/gpac_1.0.1-rev0-gd8538e8a-master_amd64.deb
RUN dpkg -i gpac_1.0.1-rev0-gd8538e8a-master_amd64.deb

ADD ./deploy/nginx-rtmp /usr/local/nginx-rtmp

WORKDIR /root/go/src/github.com/videocoin/cloud-media-server
ADD ./ ./
RUN make build

RUN cp bin/mediaserver /mediaserver-linux-amd64

WORKDIR /

ENTRYPOINT ["/mediaserver-linux-amd64"]
