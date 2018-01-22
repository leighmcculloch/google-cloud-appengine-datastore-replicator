package randid

import (
	"crypto/rand"

	"4d63.com/randstr/lib/charset"
	"4d63.com/randstr/lib/randstr"
)

var chars = charset.CharsetArray("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const length = 32

func Generate() string {
	id, err := randstr.String(rand.Reader, chars, length)
	if err != nil {
		panic(err)
	}
	return id
}
