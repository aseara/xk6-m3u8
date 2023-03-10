package m3u8

import (
	"os"
	"path/filepath"

	"github.com/aseara/xk6-m3u8/hls"
	"go.k6.io/k6/js/modules"
)

type StreamPlayer struct {
	vu       modules.VU
	recorder *hls.Recorder
	subDir   string
}

func NewPlayer(vu modules.VU) *StreamPlayer {
	rt := vu.Runtime()
	subDir := filepath.Join("tmp", rt.Get("__VU").String())
	_ = os.MkdirAll(subDir, os.ModePerm)
	return &StreamPlayer{
		vu:     vu,
		subDir: subDir,
	}
}

func (p *StreamPlayer) Set(url string) {
	p.recorder = hls.NewRecorder(url, p.subDir)
}

func (p *StreamPlayer) Record() {
	_, _ = p.recorder.Record()
}
