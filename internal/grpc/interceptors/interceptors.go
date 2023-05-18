// Package interceptors хранит interceptors для grpc.
package interceptors

import "mfuss/internal/grpc/auth"

type interceptors struct {
	a auth.Authenticator
}

// New - конструктор interceptors.
func New(a auth.Authenticator) interceptors {
	return interceptors{
		a: a,
	}
}
