package m3u8

import (
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

type (
	RootModule struct{}

	ModuleInstance struct {
		vu     modules.VU
		m      *m3u8Metrics
		player *StreamPlayer
	}
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

func New() *RootModule {
	return &RootModule{}
}

func (r RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	m, err := registerMetrics(vu)
	if err != nil {
		common.Throw(vu.Runtime(), err)
	}
	return &ModuleInstance{
		vu:     vu,
		m:      m,
		player: NewPlayer(vu, m),
	}
}

func (m ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: m.player,
	}
}
