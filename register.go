package redis

import (
	"github.com/aseara/xk6-m3u8/m3u8"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/m3u8", new(m3u8.RedisExt))
}
