package utilits

import (
	"math/rand"
	"net/url"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenetareID() string {
	rand.Seed(time.Now().UnixNano())

	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(buf)
}

func CheckURL(shortURLId string) error {
	if _, err := url.ParseRequestURI(shortURLId); err != nil {

		return err
	}
	return nil
}
