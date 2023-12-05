package connectors

import "fmt"

var (
	ErrConnectorNotFound *ConnectorNotFoundError
)

type ConnectorNotFoundError struct {
	ConnectorType string
	ConnectorName string
}

func (e *ConnectorNotFoundError) Error() string {
	return fmt.Sprintf("connector %v %v not found", e.ConnectorType, e.ConnectorName)
}
