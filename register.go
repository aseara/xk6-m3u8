package m3u8

import (
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/m3u8", New())
}
