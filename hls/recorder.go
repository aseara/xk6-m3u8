package hls

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Recorder struct {
	client        *http.Client
	dir           string
	url           string
	pulledSegment map[uint64]bool
}

func NewRecorder(url string, dir string) *Recorder {
	return &Recorder{
		url:           url,
		dir:           dir,
		client:        &http.Client{},
		pulledSegment: map[uint64]bool{},
	}
}

// Start starts a record a live-streaming
func (r *Recorder) Record() (string, error) {
	puller, d := r.pullSegment(r.url)
	d = d + rand.Float64()*2
	duration := time.NewTicker(time.Duration(d) * time.Second)

	filePath := filepath.Join(r.dir, "video.ts")
	log.Printf("Start record live streaming movie at %s...", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	for segment := range puller {
		if segment.Err != nil {
			return "", segment.Err
		}

		dc := r.downloadSegmentC(segment.Segment)

		select {
		case report := <-dc:
			if report.Err != nil {
				return "", report.Err
			}

			if _, err := file.Write(report.Data); err != nil {
				return "", err
			}
		}

		log.Println("Recorded segment ", segment.Segment.SeqId)
	}
	<-duration.C
	return filePath, nil
}

type DownloadSegmentReport struct {
	Data []byte
	Err  error
}

func (r *Recorder) downloadSegmentC(segment *Segment) chan *DownloadSegmentReport {
	c := make(chan *DownloadSegmentReport, 1)
	go func() {
		data, err := r.downloadSegment(segment)
		c <- &DownloadSegmentReport{
			Data: data,
			Err:  err,
		}
	}()

	return c
}

func (r *Recorder) downloadSegment(segment *Segment) ([]byte, error) {

	res, err := r.client.Get(segment.URI)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if segment.Key != nil {
		key, iv, err := r.getKey(segment)
		if err != nil {
			return nil, err
		}
		data, err = decryptAES128(data, key, iv)
		if err != nil {
			return nil, err
		}
	}

	for j := 0; j < len(data); j++ {
		if data[j] == syncByte {
			data = data[j:]
			break
		}
	}

	return data, nil
}

func (r *Recorder) getKey(segment *Segment) (key []byte, iv []byte, err error) {

	res, err := r.client.Get(segment.Key.URI)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != 200 {
		return nil, nil, errors.New("failed to get decryption key")
	}

	key, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	iv = []byte(segment.Key.IV)
	if len(iv) == 0 {
		iv = defaultIV(segment.SeqId)
	}
	return
}
