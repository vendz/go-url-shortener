package helper

import (
	"math/rand"
	"strings"
)

func Sug() string {
	const base58Chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	var builder strings.Builder
	base58CharsLength := len(base58Chars)

	for i := 0; i < 8; i++ {
		randomIndex := rand.Intn(base58CharsLength)
		builder.WriteByte(base58Chars[randomIndex])
	}

	return builder.String()
}
