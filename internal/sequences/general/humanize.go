package general

import (
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func HumanizedMillis(v int) time.Duration {
	if v == 0 {
		return 0
	}

	var d = v / 10
	var v1 = v - d

	var v2 = v1 + r.Intn(d)*2

	return time.Millisecond * time.Duration(v2)
}
