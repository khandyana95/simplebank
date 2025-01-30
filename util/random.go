package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(l int) string {
	var sb strings.Builder

	for i := 0; i < l; i++ {
		c := alphabet[randomInt(0, int64(len(alphabet)-1))]
		// c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return randomString(10)
}

func RandomName() string {
	return randomString(10) + " " + randomString(10)
}

func RandomMoney() float64 {
	return float64(randomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{INR, USD, CAD}

	//return currencies[randomInt(0, int64(len(currencies)-1))]
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	return randomString(10) + "@" + randomString(4) + "." + randomString(3)
}
