package m3u8

import (
	"log"
	"os"
	"time"

	"go.k6.io/k6/js/modules"
)
import "github.com/canhlinh/hlsdl"

type StreamPlayer struct {
	vu       modules.VU
	recorder *hlsdl.Recorder
}

func (player *StreamPlayer) Start(url string, dir string) error {
	player.recorder = hlsdl.NewRecorder(url, dir)
	recordedFile, err := player.recorder.Start()
	if err != nil {
		_ = os.RemoveAll(recordedFile)
		return err
	}

	log.Println("Recorded file at ", recordedFile)
	return nil
}

func (player *StreamPlayer) Check() {
	time.Sleep(5 * time.Second)
}
