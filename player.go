package m3u8

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aseara/xk6-m3u8/hls"
	"go.k6.io/k6/js/modules"
)

type StreamPlayer struct {
	vu       modules.VU
	recorder *hls.Recorder
	subDir   string
}

func (p *StreamPlayer) Start(url string, dir string) error {
	rt := p.vu.Runtime()
	p.subDir = filepath.Join(dir, rt.Get("__VU").String())
	_ = os.Mkdir(p.subDir, os.ModePerm)
	log.Println("Recorded file at ", p.subDir)
	p.recorder = hls.NewRecorder(url, p.subDir)
	go func() {
		recordedFile, err := p.recorder.Start()
		if err != nil {
			_ = os.RemoveAll(recordedFile)
			log.Println("play error", err)
		}
	}()

	return nil
}

func (p *StreamPlayer) Check() {
	time.Sleep(5 * time.Second)
}

func (p *StreamPlayer) Stop() {
	p.recorder.Stop()
	_ = os.RemoveAll(p.subDir)
}
