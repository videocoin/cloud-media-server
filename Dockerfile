FROM ubuntu

LABEL maintainer="Videocoin" description="webrtc to rtmp streamer"

RUN apt-get update
RUN apt-get install -y --no-install-recommends dh-autoreconf sudo g++-7 git ca-certificates
RUN update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-7 90

WORKDIR /go/src/github.com/videocoin

RUN git clone --recurse-submodules https://github.com/notedit/media-server-go-native.git

WORKDIR /go/src/github.com/videocoin/media-server-go-native

RUN make ssl 
RUN make srtp 
RUN make mp4
RUN make mediaserver

WORKDIR /go/src/github.com/videocoin

RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-go

RUN git clone https://github.com/notedit/media-server-go.git

WORKDIR /go/src/github.com/videocoin/media-server-go

RUN apt-get install -y \
    libgstreamer-plugins-base1.0-dev \
    libgstreamer1.0-0 \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    gstreamer1.0-tools \
    gstreamer1.0-x \
    gstreamer1.0-pulseaudio

RUN go build 

WORKDIR /go/src/github.com/videocoin/mediaserver

ADD ./ ./

RUN make build

RUN cp /go/src/github.com/videocoin/mediaserver/bin/mediaserver /mediaserver

WORKDIR /

CMD ["/mediaserver"]
