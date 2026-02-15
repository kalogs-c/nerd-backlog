package auth

import "errors"

var ErrInvalidHashFormat = errors.New("invalid hash format")
var ErrInvalidToken = errors.New("invalid token")
