package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length int) string {
	var str strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		str.WriteByte(alphabet[rand.Intn(k)])
	}

	return str.String()
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}
