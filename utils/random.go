package utils

import (
	"math/rand"
	"strings"
	"time"
)


const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init(){
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString() string{
	var s strings.Builder
	k := len(alphabet)

	for i := 0; i < 10; i++ {
		c := alphabet[rand.Intn(k)]
		s.WriteByte(c)
	}

	return s.String()
}

func RandomOwner() string {
	return RandomString()
}

func RandomMoney() int64 {
	return int64(randomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

