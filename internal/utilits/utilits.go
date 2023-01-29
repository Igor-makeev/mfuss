package utilits

import (
	"fmt"
	"math/rand"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenetareID() string {
	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(buf)
}

type URLConflict struct {
	Str string
}

func (is URLConflict) Error() string {
	return fmt.Sprintf("error:  url: %v has already been shortened", is.Str)
}

func CheckURL(shortURLId string) error {
	if _, err := url.ParseRequestURI(shortURLId); err != nil {

		return err
	}
	return nil
}
