package customErrors

import "errors"

var ErrUserAlreadyExists error = errors.New("user already exists")
