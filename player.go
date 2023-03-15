package m3u8

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/aseara/xk6-m3u8/hls"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type StreamPlayer struct {
	vu       modules.VU
	m        *m3u8Metrics
	recorder *hls.Recorder
	subDir   string
}

func NewPlayer(vu modules.VU, m *m3u8Metrics) *StreamPlayer {
	rt := vu.Runtime()
	subDir := filepath.Join("tmp", rt.Get("__VU").String())
	_ = os.MkdirAll(subDir, os.ModePerm)
	return &StreamPlayer{
		vu:     vu,
		m:      m,
		subDir: subDir,
	}
}

func (p *StreamPlayer) Set(url string) error {
	p.recorder = hls.NewRecorder(url, p.subDir)
	_, err := p.recorder.Record()
	if err != nil {
		rt := p.vu.Runtime()
		common.Throw(rt, err)
	}
	return err
}

func (p *StreamPlayer) Record() error {
	_, err := p.recorder.Record()
	_ = p.recordMetrics(err)
	return err
}

func (p *StreamPlayer) recordMetrics(err error) error {
	now := time.Now()
	state := p.vu.State()
	if state == nil {
		return errors.New("invalid state")
	}
	ctx := p.vu.Context()
	if ctx == nil {
		return errors.New("invalid state")
	}
	var failed float64
	if err != nil {
		failed = 1
	}
	metrics.PushIfNotDone(ctx, state.Samples, metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: p.m.RecordFailed,
			Tags:   p.m.TagsAndMeta.Tags,
		},
		Time:     now,
		Metadata: p.m.TagsAndMeta.Metadata,
		Value:    failed,
	})
	return nil
}
