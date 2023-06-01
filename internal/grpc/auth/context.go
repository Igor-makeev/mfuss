package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrValueIsNotUUID - значение не может быть преобразовано к типу uuid.UUID.
var ErrValueIsNotUUID = errors.New("value is not uuid.UUID")

// GetUser - функция, чтобы получить пользователя из контекста.
func GetUser(ctx context.Context) (user uuid.UUID, err error) {
	user, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrValueIsNotUUID
	}
	return
}
