package skipset

import (
	"github.com/bytedance/gg/internal/fastrand"
)

const (
	maxLevel            = 16
	p                   = 0.25
	defaultHighestLevel = 3
)

func randomLevel() int {
	level := 1
	for fastrand.Uint32n(1/p) == 0 {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}
