package webrtc

import "C"

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	msgo "github.com/notedit/media-server-go"
	"github.com/notedit/sdp"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/mediaserver/mediacore/gstrtmp"
)

func init() {
	msgo.EnableDebug(true)
	msgo.EnableLog(true)
}

var Capabilities = map[string]*sdp.Capability{
	"audio": {
		Codecs: []string{"opus"},
	},
	"video": {
		Codecs: []string{"h264"},
		Rtx:    true,
		Rtcpfbs: []*sdp.RtcpFeedback{
			{
				ID: "goog-remb",
			},
			{
				ID: "transport-cc",
			},
			{
				ID:     "ccm",
				Params: []string{"fir"},
			},
			{
				ID:     "nack",
				Params: []string{"pli"},
			},
		},
		Extensions: []string{
			"urn:3gpp:video-orientation",
			"http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01",
			"http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time",
			"urn:ietf:params:rtp-hdrext:toffse",
			"urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id",
			"urn:ietf:params:rtp-hdrext:sdes:mid",
		},
	},
}

type Streamer struct {
	logger   *logrus.Entry
	endpoint *msgo.Endpoint
	rtmpURL  string
}

func NewStreamer(ctx context.Context, host string, port int, rtmpURL string) (*Streamer, error) {
	ms := &Streamer{
		logger:   ctxlogrus.Extract(ctx).WithField("system", "webrtc-streamer"),
		endpoint: msgo.NewEndpointWithPort(host, port),
		rtmpURL:  rtmpURL,
	}

	return ms, nil
}

func (ms *Streamer) Start() error {
	return nil
}

func (ms *Streamer) Stop() error {
	ms.endpoint.Stop()
	return nil
}

func (ms *Streamer) CreateStream(sdpOffer string) (*sdp.StreamInfo, *msgo.IncomingStream, *sdp.SDPInfo, error) {
	offer, err := sdp.Parse(sdpOffer)
	if err != nil {
		return nil, nil, nil, err
	}

	transport := ms.endpoint.CreateTransport(offer, nil)
	transport.SetRemoteProperties(offer.GetMedia("audio"), offer.GetMedia("video"))

	answer := offer.Answer(
		transport.GetLocalICEInfo(),
		transport.GetLocalDTLSInfo(),
		ms.endpoint.GetLocalCandidates(),
		Capabilities,
	)

	transport.SetLocalProperties(answer.GetMedia("audio"), answer.GetMedia("video"))

	stream := offer.GetFirstStream()
	incomingStream := transport.CreateIncomingStream(stream)

	return stream, incomingStream, answer, nil
}

func (ms *Streamer) StartStreaming(
	streamID string,
	stream *sdp.StreamInfo,
	incomingStream *msgo.IncomingStream,
) error {
	logger := ms.logger.WithField("stream_id", streamID)

	go func() {
		refresher := msgo.NewRefresher(2000)
		refresher.AddStream(incomingStream)
		// defer refresher.Stop()

		rtmpURL := fmt.Sprintf("%s/%s", ms.rtmpURL, streamID)
		logger.Infof("re-streaming to rtmp url %s", rtmpURL)

		pipeline := gstrtmp.CreatePipeline(rtmpURL)
		pipeline.Start()

		if len(incomingStream.GetVideoTracks()) > 0 {
			videoTrack := incomingStream.GetVideoTracks()[0]
			videoTrack.OnMediaFrame(func(frame []byte, timestamp uint) {
				if len(frame) <= 4 {
					return
				}
				pipeline.PushVideo(frame)
			})
		}

		if len(incomingStream.GetAudioTracks()) > 0 {
			audioTrack := incomingStream.GetAudioTracks()[0]
			audioTrack.OnMediaFrame(func(frame []byte, timestamp uint) {
				pipeline.PushAudio(frame)
			})
		}
	}()

	return nil
}
