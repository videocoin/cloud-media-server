package mediacore

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"cloud.google.com/go/storage"
)

const (
	noCache          = "no-cache"
	MimeTypeMP4      = "video/mp4"
	MimeTypeMpegDash = "application/dash+xml"
)

func MuxToMp4(streamID string, inputURL string) (string, error) {
	outputPath := fmt.Sprintf("/tmp/%s.mp4", streamID)

	execArgs := []string{
		"-i",
		inputURL,
		"-c:v",
		"copy",
		"-c:a",
		"copy",
		"-map",
		"0:v",
		"-map",
		"'0:a?'",
		"-f",
		"mp4",
		outputPath,
	}
	cmd := exec.Command("ffmpeg", execArgs...)

	fmt.Println("ffmpeg " + strings.Join(execArgs, " "))

	cmdOut, err := cmd.Output()
	if err != nil {
		fmtErr := err
		if len(string(cmdOut)) > 0 {
			fmtErr = fmt.Errorf("%s: %s", err.Error(), string(cmdOut))
		}
		return "", fmtErr
	}

	return outputPath, nil
}

func UploadFileToBucket(ctx context.Context, bh *storage.BucketHandle, streamID string, filename string, ct string, src *os.File) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/%s", streamID, filename)

	obj := bh.Object(objectName)
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

	return obj, attrs, err
}

func UploadFilepathToBucket(ctx context.Context, bh *storage.BucketHandle, streamID string, filename string, ct string, src string) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/%s", streamID, filename)

	obj := bh.Object(objectName)
	w := obj.NewWriter(ctx)
	w.CacheControl = noCache
	w.ContentType = ct

	f, err := os.Open(src)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
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

	return obj, attrs, err
}
