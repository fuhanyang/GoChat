package random

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}
func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
