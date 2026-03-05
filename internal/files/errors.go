package files

import "errors"

var (
	ErrObjectAlreadyExists = errors.New("Folder already exists")
	ErrExpectedDir         = errors.New("expected a directory, found a file")
)
