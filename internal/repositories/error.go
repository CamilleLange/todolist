package repositories

import "fmt"

var (
	ErrFeatureNotImplemented = fmt.Errorf("feature not implemented")
	ErrNoRowAffected         = fmt.Errorf("no row affected")

	ErrDAOTypeNotFound *DAOTypeNotFoundError
	ErrNoDataFound     *NoDataFoundError
)

type DAOTypeNotFoundError struct {
	Type string
}

func (e *DAOTypeNotFoundError) Error() string {
	return fmt.Sprintf("type %v not found", e.Type)
}

type FailToBuildDAOError struct {
	Type string
	Err  error
}

type NoDataFoundError struct{}

func (e *NoDataFoundError) Error() string {
	return "no data found"
}

type InvalidContextError struct {
	Key string
}

func (e *InvalidContextError) Error() string {
	return "requested object " + e.Key + " is not present or is not castable to destination struct"
}
