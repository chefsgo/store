package file

import "errors"

const (
	NAME = "STORE"
)

var (
	errInvalidStoreConnection = errors.New("Invalid file connection.")
)
