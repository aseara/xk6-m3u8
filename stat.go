package m3u8

import (
	"errors"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type m3u8Metrics struct {
	RecordFailed *metrics.Metric
	TagsAndMeta  *metrics.TagsAndMeta
}

// registerMetrics registers the metrics for the mqtt module in the metrics registry
func registerMetrics(vu modules.VU) (*m3u8Metrics, error) {
	var err error
	m := &m3u8Metrics{}
	env := vu.InitEnv()
	if env == nil {
		return m, errors.New("missing env")
	}
	registry := env.Registry
	if registry == nil {
		return m, errors.New("missing registry")
	}
	m.RecordFailed, err = registry.NewMetric("m3u8_record_failed", metrics.Rate)
	if err != nil {
		return m, err
	}
	m.TagsAndMeta = &metrics.TagsAndMeta{
		Tags: registry.RootTagSet(),
	}
	return m, err
}
