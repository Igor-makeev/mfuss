package utilits

import (
	"math/rand"
	"net/url"
	"time"
)

// строка из которой рандомно выбираются числа
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// рандомим
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// ссылка проверяющая валидность урла
func CheckURL(shortURLId string) error {
	if _, err := url.ParseRequestURI(shortURLId); err != nil {

		return err
	}
	return nil
}

// генератор ID
func GenetareID() string {

	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = letterBytes[r.Intn(len(letterBytes))]
	}

	return string(buf)
}
