package repository

import "errors"

var (
	ErrNoRows          = errors.New("no rows returned")
	ErrQueryExec       = errors.New("failed to execute query")
	ErrRowScan         = errors.New("error scanning rows")
	ErrNilEntity       = errors.New("entity is nil")
	ErrDependentEntity = errors.New("entity is dependent")
)
