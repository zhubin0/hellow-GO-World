package oss

import "errors"

var ErrNotFound = errors.New("object not found")

func IsNotFound(err error) bool {
	return err == ErrNotFound
}
