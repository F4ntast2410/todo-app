package customErrors

import "errors"

var ErrInccorectPassword error = errors.New("wrong old password")
