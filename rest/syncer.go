package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/grafov/m3u8"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	pStreamsV1 "github.com/videocoin/cloud-api/streams/private/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-media-server/datastore"
)

const (
	noCache        = "no-cache"
	mimeTypeM3U8   = "application/x-mpegURL"
	mimeTypeMPEGTS = "video/mp2t"
	mimeTypeMP4    = "video/mp4"
)

func (s *Server) sync(c echo.Context) error {
	path := c.FormValue("path")
	ct := strings.ToLower(c.FormValue("ct"))
	vod := c.FormValue("vod")
	last := c.FormValue("last")
	durationStr := c.FormValue("duration")
	duration, _ := strconv.ParseFloat(durationStr, 64)
	isLast := last == "y"
	isVOD := vod == "y"

	logger := s.logger.WithFields(logrus.Fields{
		"path":     path,
		"ct":       ct,
		"vod":      vod,
		"duration": durationStr,
		"is_last":  isLast,
		"is_vod":   isVOD,
	})

	logger.Info("syncing")

	file, err := c.FormFile("file")
	if err != nil {
		logger.Errorf("failed to get file: %s", err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		logger.Errorf("failed to open file: %s", err)
		return err
	}
	defer src.Close()

	streamID, segmentNum, err := parseReqPath(path)
	if err != nil {
		e := fmt.Errorf("failed to parse request path: %s", err)
		logger.Error(e)
		return e
	}

	emptyCtx := context.Background()
	_, _, err = s.uploadSegment(emptyCtx, streamID, segmentNum, ct, src)
	if err != nil {
		e := fmt.Errorf("failed to upload segment: %s", err)
		logger.Error(e)
		return e
	}

	err = s.ds.AddSegment(streamID, segmentNum, duration)
	if err != nil {
		e := fmt.Errorf("failed to add segment: %s", err.Error())
		logger.Error(e)
		return e
	}

	segments, err := s.ds.GetSegments(streamID)
	if err != nil {
		e := fmt.Errorf("failed to get segments: %s", err.Error())
		logger.Error(e)
		return e
	}

	if !isVOD {
		logger.Info("generating and uploading live master playlist")
		_, _, err = s.generateAndUploadLiveMasterPlaylist(emptyCtx, streamID, segments, ct, isLast)
		if err != nil {
			e := fmt.Errorf("failed to generate live master playlist: %s", err.Error())
			s.logger.Error(e)
			return e
		}

		if segmentNum == 1 {
			logger.Info("updating stream status as ready")
			err = s.eb.EmitUpdateStreamStatus(emptyCtx, streamID, streamsv1.StreamStatusReady)
			if err != nil {
				logger.Errorf("failed to update stream status: %s", err)
			}
		}
	} else {
		if len(segments) > 0 {
			logger.Info("generating and uploading vod master playlist")
			_, _, err = s.generateAndUploadVODMasterPlaylist(emptyCtx, streamID, segments, ct)
			if err != nil {
				e := fmt.Errorf("failed to generate vod master playlist: %s", err.Error())
				s.logger.Error(e)
				return e
			}

			streamResp, err := s.sc.Streams.Get(emptyCtx, &pStreamsV1.StreamRequest{Id: streamID})
			if err != nil {
				e := fmt.Errorf("failed to get stream: %s", err.Error())
				s.logger.Error(e)
				return e
			}

			if streamResp.OutputType == streamsv1.OutputTypeDash ||
				streamResp.OutputType == streamsv1.OutputTypeDashWithDRM {

				logger.Info("generating and uploading mpd manifest")
				_, _, err = s.generateAndUploadMPD(emptyCtx, streamID, segments, ct)
				if err != nil {
					e := fmt.Errorf("failed to generate mpd manifest: %s", err.Error())
					s.logger.Error(e)
					return e
				}

			}
		}
	}

	return c.NoContent(http.StatusAccepted)
}

func (s *Server) uploadSegment(ctx context.Context, streamID string, segmentNum int, ct string, src multipart.File) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	ext := ".ts"
	if ct == mimeTypeMP4 {
		ext = ".mp4"
	}
	objectName := fmt.Sprintf("%s/%d%s", streamID, segmentNum, ext)

	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"segment_num": segmentNum,
		"bucket":      s.bucket,
		"object_name": objectName,
	})

	logger.Info("uploading segment")

	obj := s.bh.Object(objectName)
	w := obj.NewWriter(ctx)
	w.CacheControl = noCache
	w.ContentType = ct

	if _, err := io.Copy(w, src); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return obj, attrs, err
	}

	logger.Info("segment has been uploaded successfully")

	return obj, attrs, err
}

func (s *Server) generateAndUploadLiveMasterPlaylist(
	ctx context.Context,
	streamID string,
	segments []*datastore.Segment,
	ct string,
	last bool,
) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/index.m3u8", streamID)
	tmpObjectName := fmt.Sprintf("%s/_index.m3u8", streamID)

	ext := ".ts"
	if ct == mimeTypeMP4 {
		ext = ".mp4"
	}

	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"bucket":      s.bucket,
		"object_name": objectName,
		"last":        last,
		"ct":          ct,
	})

	logger.Info("generating live master playlist")

	segmentsCount := len(segments)
	p, err := m3u8.NewMediaPlaylist(uint(segmentsCount), uint(segmentsCount))
	if err != nil {
		return nil, nil, err
	}

	if last {
		p.MediaType = m3u8.VOD
	}

	for _, segment := range segments {
		err := p.Append(fmt.Sprintf("%d%s", segment.Num, ext), segment.Duration, "")
		if err != nil {
			return nil, nil, err
		}
	}

	if last {
		p.Close()
	}

	data := p.Encode().Bytes()

	logger.Info("uploading live master playlist")

	obj := s.bh.Object(tmpObjectName)

	w := obj.NewWriter(ctx)
	w.CacheControl = noCache
	w.ContentType = mimeTypeM3U8

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	dst := s.bh.Object(objectName)

	copier := dst.CopierFrom(obj)
	copier.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	copier.ContentType = mimeTypeM3U8
	copier.CacheControl = "private, max-age=0, no-transform"

	_, err = copier.Run(context.Background())
	if err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return obj, attrs, err
	}

	logger.Info("live master playlist has been uploaded successfully")

	return obj, attrs, err
}

func (s *Server) generateAndUploadVODMasterPlaylist(
	ctx context.Context,
	streamID string,
	segments []*datastore.Segment,
	ct string,
) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/index.m3u8", streamID)
	tmpObjectName := fmt.Sprintf("%s/_index.m3u8", streamID)

	ext := ".ts"
	if ct == mimeTypeMP4 {
		ext = ".mp4"
	}

	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"bucket":      s.bucket,
		"object_name": objectName,
		"ct":          ct,
	})

	logger.Info("generating vod master playlist")

	segmentsCount := len(segments)
	p, err := m3u8.NewMediaPlaylist(uint(segmentsCount), uint(segmentsCount))
	if err != nil {
		return nil, nil, err
	}

	p.MediaType = m3u8.VOD

	for _, segment := range segments {
		err := p.Append(fmt.Sprintf("%d%s", segment.Num, ext), segment.Duration, "")
		if err != nil {
			return nil, nil, err
		}
	}

	p.Close()

	data := p.Encode().Bytes()

	logger.Info("uploading vod master playlist")

	obj := s.bh.Object(tmpObjectName)

	w := obj.NewWriter(ctx)
	w.CacheControl = noCache
	w.ContentType = mimeTypeM3U8

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	dst := s.bh.Object(objectName)

	copier := dst.CopierFrom(obj)
	copier.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	copier.ContentType = mimeTypeM3U8
	copier.CacheControl = "private, max-age=0, no-transform"

	_, err = copier.Run(context.Background())
	if err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return obj, attrs, err
	}

	logger.Info("vod master playlist has been uploaded successfully")

	return obj, attrs, err
}

func (s *Server) generateAndUploadMPD(
	ctx context.Context,
	streamID string,
	segments []*datastore.Segment,
	ct string,
) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/index.m3u8", streamID)
	tmpObjectName := fmt.Sprintf("%s/_index.mpd", streamID)

	ext := ".ts"
	if ct == mimeTypeMP4 {
		ext = ".mp4"
	}

	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"bucket":      s.bucket,
		"object_name": objectName,
		"ct":          ct,
	})

	logger.Info("generating mpd manifest")

	segmentsCount := len(segments)
	p, err := m3u8.NewMediaPlaylist(uint(segmentsCount), uint(segmentsCount))
	if err != nil {
		return nil, nil, err
	}

	p.MediaType = m3u8.VOD

	for _, segment := range segments {
		err := p.Append(fmt.Sprintf("%d%s", segment.Num, ext), segment.Duration, "")
		if err != nil {
			return nil, nil, err
		}
	}

	p.Close()

	data := p.Encode().Bytes()

	logger.Info("uploading mpd manifest")

	obj := s.bh.Object(tmpObjectName)

	w := obj.NewWriter(ctx)
	w.CacheControl = noCache
	w.ContentType = mimeTypeM3U8

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	dst := s.bh.Object(objectName)

	copier := dst.CopierFrom(obj)
	copier.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	copier.ContentType = mimeTypeM3U8
	copier.CacheControl = "private, max-age=0, no-transform"

	_, err = copier.Run(context.Background())
	if err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return obj, attrs, err
	}

	logger.Info("mpd manifest has been uploaded successfully")

	return obj, attrs, err
}

func parseReqPath(path string) (string, int, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", 0, errors.New("wrong request path format")
	}

	streamID := parts[0]

	sparts := strings.Split(parts[1], ".")
	if len(sparts) != 2 {
		return "", 0, errors.New("wrong segment format")
	}

	segmentNum, err := strconv.Atoi(sparts[0])
	if err != nil {
		return "", 0, err
	}

	return streamID, segmentNum, nil
}
