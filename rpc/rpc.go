package rpc

import (
	"context"

	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/mediaserver/v1"
	"github.com/videocoin/cloud-api/rpc"
)

func (s *RpcServer) CreateWebRTCStream(ctx context.Context, req *v1.StreamRequest) (*v1.WebRTCStreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.StreamId)
	span.SetTag("sdp", req.Sdp)

	userID := 0
	// userID, _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	span.SetTag("user_id", userID)

	logger := s.logger.
		WithField("stream_id", req.StreamId).
		WithField("user_id", userID)

	if verr := s.validator.validate(req); verr != nil {
		logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	sdpStream, inStream, sdpAnswer, err := s.ms.CreateWebRTCStream(req.Sdp)
	if err != nil {
		return nil, err
	}

	logger.Debugf("sdpStream %+v", sdpStream)
	logger.Debugf("incomingStream %+v", inStream)

	err = s.ms.StartWebRTCStreaming(req.StreamId, sdpStream, inStream)
	if err != nil {
		return nil, err
	}

	resp := &v1.WebRTCStreamResponse{
		Sdp: sdpAnswer.String(),
	}

	return resp, nil
}
