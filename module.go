package m3u8

import (
	"go.k6.io/k6/js/modules"
)

type (
	RootModule struct{}

	ModuleInstance struct {
		vu     modules.VU
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
	return &ModuleInstance{
		vu:     vu,
		player: NewPlayer(vu),
	}
}

func (m ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: m.player,
	}
}
