package rand

import (
	"math/rand"
	"time"
)

// Random functions
type Random struct{}

// String return random string
func (h *Random) String(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Int return random int from string slice
func (h *Random) Int(max int) int {
	return rand.Intn(max)
}

// Float return random float64 from string slice
func (h *Random) Float(max int) float64 {
	return float64(rand.Intn(max)) + rand.Float64()
}

// Bool return random bool
func (h *Random) Bool() bool {
	boolList := []bool{true, false}
	randomIndex := rand.Intn(len(boolList))
	return boolList[randomIndex]
}

// DateTime return random time
func (h *Random) DateTime() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// StringFromSlice return random string from string slice
func (h *Random) StringFromSlice(in []string) string {
	randomIndex := rand.Intn(len(in))
	return in[randomIndex]
}

// IntFromSlice return random int from string slice
func (h *Random) IntFromSlice(in []int64) int64 {
	randomIndex := rand.Intn(len(in))
	return in[randomIndex]
}
