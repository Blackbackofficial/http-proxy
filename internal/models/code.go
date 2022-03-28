package models

type StatusCode int

const (
	Okey StatusCode = iota
	Created
	NotFound
	InternalError
	Conflict
	BadRequest
)
