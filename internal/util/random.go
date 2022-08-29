// Package util contains utility functions.
package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

// RandomString returns a random string of length n.
// The string is generated from the alphabet defined in alphabet.
// If n is less than 1, an empty string is returned.
// If n is greater than the length of the alphabet, an error is returned.
// If n is not an integer, an error is returned.
func RandomString(length int) string {
	var str strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		str.WriteByte(alphabet[rand.Intn(k)])
	}

	return str.String()
}

// RandomInt returns a random integer between min and max.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}
