package rpc

import (
	"context"
	"os"

	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/mediaserver/v1"
	"github.com/videocoin/cloud-api/rpc"
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-media-server/mediacore"
)

func (s *Server) CreateWebRTCStream(ctx context.Context, req *v1.StreamRequest) (*v1.WebRTCStreamResponse, error) {
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

func (s *Server) Mux(ctx context.Context, req *v1.MuxRequest) (*v1.MuxResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.StreamId)
	span.SetTag("bucket", s.bucket)

	logger := s.logger.
		WithField("stream_id", req.StreamId).
		WithField("bucket", s.bucket)

	if req.StreamId == "" {
		return nil, rpc.ErrRpcBadRequest
	}

	go func() {
		logger.Info("getting stream")

		streamReq := &pstreamsv1.StreamRequest{Id: req.StreamId}
		streamResp, err := s.streams.Get(context.Background(), streamReq)
		if err != nil {
			logger.WithError(err).Error("failed to get stream")
			return
		}

		span.SetTag("input_url", streamResp.OutputURL)
		logger = logger.WithField("input_url", streamResp.OutputURL)

		logger.Info("muxing")

		outPath, err := mediacore.MuxToMp4(req.StreamId, streamResp.OutputURL)
		if err != nil {
			logger.WithError(err).Error("failed to mux to file")
			return
		}

		logger.WithField("output", outPath).Info("mux has been completed")

		f, err := os.Open(outPath)
		if err != nil {
			logger.WithError(err).Error("failed to open file")
			return
		}
		defer func() {
			err = f.Close()
			if err != nil {
				logger.WithError(err).Error("failed to close mux file")
				return
			}

			err = os.Remove(outPath)
			if err != nil {
				logger.WithError(err).Error("failed to remove mux file")
				return
			}
		}()

		logger.Info("uploading mux file to bucket")

		emptyCtx := context.Background()
		_, _, err = mediacore.UploadFileToBucket(
			emptyCtx,
			s.bh,
			req.StreamId,
			"index.mp4",
			mediacore.MimeTypeMP4,
			f,
		)
		if err != nil {
			logger.WithError(err).Error("failed to upload file to bucket")
			return
		}

		logger.Info("upload mux file has been completed")

		_, err = s.streams.Complete(emptyCtx, &pstreamsv1.StreamRequest{Id: req.StreamId})
		if err != nil {
			logger.WithError(err).Error("failed to complete stream")
			return
		}
	}()

	return &v1.MuxResponse{}, nil
}
