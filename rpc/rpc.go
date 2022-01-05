package rpc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-media-server/mediacore"

	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/mediaserver/v1"
	"github.com/videocoin/cloud-api/rpc"
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

	sdpStream, inStream, sdpAnswer, err := s.webRtcStreamer.CreateStream(req.Sdp)
	if err != nil {
		return nil, err
	}

	logger.Debugf("sdpStream %+v", sdpStream)
	logger.Debugf("incomingStream %+v", inStream)

	err = s.webRtcStreamer.StartStreaming(req.StreamId, sdpStream, inStream)
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
	span.SetTag("input_url", req.InputUrl)
	span.SetTag("bucket", s.bucket)

	emptyResp := &v1.MuxResponse{}

	logger := s.logger.
		WithField("stream_id", req.StreamId).
		WithField("input_url", req.InputUrl).
		WithField("bucket", s.bucket)

	if req.StreamId == "" {
		return nil, rpc.ErrRpcBadRequest
	}

	logger.Info("getting stream")

	streamReq := &pstreamsv1.StreamRequest{Id: req.StreamId}
	streamResp, err := s.sc.Streams.Get(context.Background(), streamReq)
	if err != nil {
		logger.WithError(err).Error("failed to get stream")
		return emptyResp, err
	}

	if streamResp.OutputType == streamsv1.OutputTypeHLS {
		return emptyResp, nil
	}

	if streamResp.OutputType == streamsv1.OutputTypeFile {
		err = s.muxToFile(context.Background(), req, streamResp)
		return emptyResp, nil
	}

	if streamResp.OutputType == streamsv1.OutputTypeDash ||
		streamResp.OutputType == streamsv1.OutputTypeDashWithDRM {

		useDrm := streamResp.OutputType == streamsv1.OutputTypeDashWithDRM
		err = s.muxToMpegDash(context.Background(), req, streamResp, useDrm)

		return emptyResp, nil
	}

	_, err = s.sc.Streams.Complete(context.Background(), &pstreamsv1.StreamRequest{Id: req.StreamId})
	if err != nil {
		logger.WithError(err).Error("failed to complete stream")
		return emptyResp, err
	}

	return emptyResp, nil
}

func (s *Server) muxToFile(ctx context.Context, req *v1.MuxRequest, stream *pstreamsv1.StreamResponse) error {
	inputURL := strings.Replace(stream.OutputURL, "/index.mp4", "/index.m3u8", -1)

	logger := s.logger.
		WithField("stream_id", req.StreamId).
		WithField("bucket", s.bucket).
		WithField("input_url", inputURL)

	logger.Info("muxing to mp4 file")

	outPath, err := mediacore.MuxToMp4(req.StreamId, inputURL)
	if err != nil {
		logger.WithError(err).Error("failed to mux to file")
		return err
	}

	logger.WithField("output", outPath).Info("mux has been completed")

	f, err := os.Open(outPath)
	if err != nil {
		logger.WithError(err).Error("failed to open file")
		return err
	}

	defer func() {
		_ = f.Close()
		_ = os.Remove(outPath)
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
		return err
	}

	logger.Info("upload mux file has been completed")

	return nil
}

func (s *Server) muxToMpegDash(ctx context.Context, req *v1.MuxRequest, stream *pstreamsv1.StreamResponse, useDrm bool) error {

	logger := s.logger.
		WithField("stream_id", req.StreamId).
		WithField("bucket", s.bucket).
		WithField("input", stream.OutputURL)

	logger.Info("generating mpeg dash")

	fpMp4, err := mediacore.MuxToMp4(req.StreamId, stream.OutputURL)
	if err != nil {
		logger.WithError(err).Error("failed to mux hls to mp4 file")
		return err
	}

	fnEcryptedMp4 := stream.ID + "_encrypted.mp4"
	fpEncryptedMp4 := path.Join("/tmp", fnEcryptedMp4)

	defer func() {
		_ = os.Remove(fpMp4)
	}()

	if useDrm {
		if stream.DrmXml == "" {
			return errors.New("empty drm xml")
		}

		fnDrmXml := "drm_" + stream.ID + ".xml"
		fpDrmXml := path.Join("/tmp", fnDrmXml)
		err = ioutil.WriteFile(fpDrmXml, []byte(stream.DrmXml), 0644)
		if err != nil {
			return fmt.Errorf("failed to write drm xml: %s", err)
		}

		_, err := mediacore.MP4BoxCryptExec(fpDrmXml, fpMp4, fpEncryptedMp4)
		if err != nil {
			return fmt.Errorf("failed to crypt mp4: %s", err)
		}
		os.Remove(fpMp4)
		os.Rename(fpEncryptedMp4, fpMp4)
	}

	fnMpd := stream.ID + ".mpd"
	fpMpd := path.Join("/tmp", fnMpd)
	_, err = mediacore.MP4BoxDashExec(fpMp4, "", fpMpd)
	if err != nil {
		logger.WithError(err).Error("failed to generate mpeg dash")
		return err
	}

	defer func() {
		_ = os.Remove(fpMpd)
	}()

	logger.Info("uploading mpd manisfest")

	_, _, err = mediacore.UploadFilepathToBucket(
		ctx,
		s.bh,
		req.StreamId,
		"index.mpd",
		mediacore.MimeTypeMpegDash,
		fpMpd,
	)
	if err != nil {
		logger.WithError(err).Error("failed to upload mpd manifest to bucket")
		return err
	}

	logger.Info("uploading mpd init segment")

	fnInitSegment := fmt.Sprintf("%sinit.mp4", stream.ID)
	fpInitSegment := strings.ReplaceAll(
		fpMp4,
		fmt.Sprintf("%s.mp4", stream.ID),
		fnInitSegment,
	)
	_, _, err = mediacore.UploadFilepathToBucket(
		ctx,
		s.bh,
		req.StreamId,
		fnInitSegment,
		mediacore.MimeTypeMP4,
		fpInitSegment,
	)
	if err != nil {
		logger.WithError(err).Error("failed to upload mpd init segment to bucket")
		return err
	}

	defer func() {
		os.Remove(fpInitSegment)
	}()

	return nil
}
