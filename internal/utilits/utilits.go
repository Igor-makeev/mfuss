package utilits

import (
	"math/rand"
	"net/url"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func CheckURL(shortURLId string) error {
	if _, err := url.ParseRequestURI(shortURLId); err != nil {

		return err
	}
	return nil
}

func GenetareID() string {

	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = letterBytes[r.Intn(len(letterBytes))]
	}

	return string(buf)
}
