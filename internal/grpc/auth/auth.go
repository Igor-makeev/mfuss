// Package authenticator хранит реализацию аутентификации пользователя.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"

	"mfuss/configs"
)

// ErrUnauthorized - Пользователь не авторизован.
var ErrUnauthorized = errors.New("unauthorized")

// Authenticator - структура, которая содержит функции для аутентификации.
type Authenticator struct {
	cfg *configs.Config
}

// New - конструктор для Authenticator.
func New(cfg *configs.Config) Authenticator {
	return Authenticator{
		cfg: cfg,
	}
}

// Load - функция, которая проверяет подпись строки и достает пользователя.
func (a Authenticator) Load(s string) (user uuid.UUID, err error) {
	payload, err := base64.StdEncoding.DecodeString(s)
	if err != nil || len(payload) < 16 {
		return uuid.Nil, ErrUnauthorized
	}

	h := hmac.New(sha256.New, a.cfg.CookieKey)
	h.Write(payload[:16])
	sign := h.Sum(nil)

	if !hmac.Equal(sign, payload[16:]) {
		return uuid.Nil, ErrUnauthorized
	}

	user, err = uuid.FromBytes(payload[:16])
	if err != nil {
		return uuid.Nil, ErrUnauthorized
	}

	return user, nil
}

// Gen - функция, которая генерирует нового пользователя.
func (a Authenticator) Gen() (user uuid.UUID, signed string) {
	user = uuid.New()

	b, _ := user.MarshalBinary()

	h := hmac.New(sha256.New, a.cfg.CookieKey)
	h.Write(b)
	sign := h.Sum(nil)

	signed = base64.StdEncoding.EncodeToString(append(b, sign...))

	return user, signed
}
