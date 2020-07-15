package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

type Datastore interface {
	Start() error
	Stop() error
	CreateFileMeta(context.Context, *FileMeta) error
	GetFileMeta(context.Context, string) (*FileMeta, error)
	AddSegment(string, int, float64) error
	GetSegments(string) ([]*Segment, error)
	GetMaxSegment(string) (int, error)
}

type datastore struct {
	logger *logrus.Entry
	cli    *redis.Client
}

func NewDatastore(ctx context.Context, uri string) (Datastore, error) {
	d := &datastore{
		logger: ctxlogrus.Extract(ctx).WithField("system", "datastore"),
	}

	opts, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	opts.MaxRetries = 3
	opts.PoolSize = 10

	cli := redis.NewClient(opts)
	if err != nil {
		return nil, err
	}

	d.cli = cli

	return d, nil
}

func (ds *datastore) Start() error {
	return ds.cli.Ping().Err()
}

func (ds *datastore) Stop() error {
	return ds.cli.Close()
}

func (ds *datastore) CreateFileMeta(ctx context.Context, meta *FileMeta) error {
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	err = ds.cli.Set(meta.ID, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ds *datastore) GetFileMeta(ctx context.Context, id string) (*FileMeta, error) {
	data, err := ds.cli.Get(id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	meta := new(FileMeta)
	err = json.Unmarshal(data, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (ds *datastore) AddSegment(streamID string, segmentNum int, duration float64) error {
	k := makePlaylistSegmentsKey(streamID)
	segment := &Segment{Num: segmentNum, Duration: duration}
	segmentJSON, err := json.Marshal(segment)
	if err != nil {
		return err
	}
	err = ds.cli.SAdd(k, string(segmentJSON)).Err()
	return err
}

func (ds *datastore) GetSegments(streamID string) ([]*Segment, error) {
	k := makePlaylistSegmentsKey(streamID)
	segments := []*Segment{}
	segmentsJSON, err := ds.cli.SMembers(k).Result()
	for _, segmentJSON := range segmentsJSON {
		segment := &Segment{}
		err := json.Unmarshal([]byte(segmentJSON), segment)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	sort.Sort(ByNum(segments))

	return segments, err
}

func (ds *datastore) GetMaxSegment(streamID string) (int, error) {
	k := makePlaylistSegmentsKey(streamID)
	maxNum := 0
	segments, err := ds.cli.SMembers(k).Result()
	for _, segmentStr := range segments {
		num, err := strconv.Atoi(segmentStr)
		if err == nil {
			if maxNum < num {
				maxNum = num
			}
		}
	}
	return maxNum, err
}

func makePlaylistSegmentsKey(streamID string) string {
	return fmt.Sprintf("playlists/%s/segments", streamID)
}
