package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func NewToken(n int) string {

	buff := make([]byte, n*3/4)
	if _, err := rand.Read(buff); err != nil {
		panic(err)
	}

	return Truncate(base64.RawURLEncoding.EncodeToString(buff), n)
}

func Truncate(val string, max int) string {

	if len(val) <= max {
		return val
	}

	return val[:max]
}
