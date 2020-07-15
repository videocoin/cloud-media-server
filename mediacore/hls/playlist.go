package hls

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/grafov/m3u8"
)

func ExtractSegments(path string) ([]*m3u8.MediaSegment, error) {
	segments := []*m3u8.MediaSegment{}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open playlist: %s", err)
	}

	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		return nil, fmt.Errorf("failed to decode playlist: %s", err)
	}

	switch listType {
	case m3u8.MASTER:
		return nil, errors.New("failed to playlist type")
	case m3u8.MEDIA:
		pl := p.(*m3u8.MediaPlaylist)
		for _, s := range pl.Segments {
			if s != nil {
				segments = append(segments, s)
			}
		}
	}

	return segments, nil
}
