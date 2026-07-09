package customErrors

import "errors"

var ErrCacheValueNotExists error = errors.New("code not found")
