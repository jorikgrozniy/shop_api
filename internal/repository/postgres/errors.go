package postgres

import "errors"

var (
	ErrQueryExec = errors.New("failed to execute query")
	ErrRowScan   = errors.New("error scanning rows")
	ErrNilEntity = errors.New("entity is nil")
)
