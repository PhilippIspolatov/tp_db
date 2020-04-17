package tools

import "errors"

var (
	BadRequest = errors.New("bad request")
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)
