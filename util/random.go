package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var r *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	r = rand.New(rand.NewSource(seed))

}

// RandInt return pseudo-random element from the half open [min, max)
func RandInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// RandString return string of size n containing elements of [a-z]
func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandOwner() string {
	return RandString(6)
}

func RandMoney() int64 {
	return RandInt(0, 1_000)
}

func RandCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}

	n := len(currencies)

	return currencies[r.Intn(n)]
}
