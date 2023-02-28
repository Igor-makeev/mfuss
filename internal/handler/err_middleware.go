package handler

import "errors"

var (
	ErrNoDataArray      = errors.New("data array not found")
	ErrInvalidDataArray = errors.New("data array is of ivalid type")
	ErrNoUserID         = errors.New("user id not found")
	ErrInvalidUserID    = errors.New("user id is of ivalid type")
)
