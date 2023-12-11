package sqldb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
	_ "github.com/lib/pq"              // postgres
)

const (
	tickInterval = 50 * time.Millisecond
)

// Conf represents the configuration structure for opening a connection to the
// database. It includes fields for specifying the database driver and Data
// Source Name (DSN) for connecting to the database.
type Conf struct {
	Driver string `yaml:"driver"` // Database driver (e.g., "mysql", "postgres", etc.).
	DSN    string `yaml:"dsn"`    // Data Source Name (DSN) for connecting to the database.
}

// Connector is a struct representing a database connector. It embeds a *sql.DB
// to provide the functionalities of a SQL database connection.
type Connector struct {
	*sql.DB
}

// TryConnection attempts to establish a connection to the SQL database by
// repeatedly pinging it within a specified time duration. It takes an integer
// parameter 't' representing the maximum time duration in seconds for attempting
// the connection. The function returns an error if the connection cannot be
// established within the specified time or if an error occurs during the connection
// attempt.
func (con *Connector) TryConnection(t int) error {
	// Create a ticker to ping the database at regular intervals.
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	// Set a timeout for the overall connection attempt.
	timeout := time.After(time.Duration(t) * time.Second)

	// Loop until the database connection is successful or the timeout occurs.
	for {
		select {
		case <-ticker.C:
			// Attempt to ping the database.
			err := con.Ping()
			if err == nil {
				return nil // Successful connection.
			}

		case <-timeout:
			// Return an error if the timeout is reached without a successful connection.
			return fmt.Errorf("can't ping SqlDB: timeout after %d s", t)
		}
	}
}

// Commit commits the provided SQL transaction. It takes a *sql.Tx parameter and
// attempts to commit the transaction. If the commit is successful, the function
// returns nil; otherwise, it returns an error with relevant information.
func (con *Connector) Commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("fail to commit transaction: %w", err)
	}
	return nil
}

// Exec executes the provided SQL query within the specified transaction. It
// takes a *sql.Tx, a query string, and optional query arguments. The function
// prepares the query, executes the prepared statement with the given arguments,
// and returns the sql.Result. If any error occurs during preparation or execution,
// the function returns an error with relevant information.
func (con *Connector) Exec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	var res sql.Result

	// Prepare the SQL query.
	stm, err := tx.Prepare(query)
	if err != nil {
		return res, fmt.Errorf("can't prepare query: %w", err)
	}

	// Execute the prepared statement with the provided arguments.
	res, err = stm.Exec(args...)
	if err != nil {
		return res, fmt.Errorf("error when executing prepared statement: %w", err)
	}

	return res, nil
}

// ExecQueryRow executes the provided SQL query within the specified transaction
// and returns a *sql.Row representing the result. It takes a *sql.Tx, a query
// string, and optional query arguments. The function prepares the query, creates
// a *sql.Row with the prepared statement and the given arguments, and returns it.
// If any error occurs during preparation, the function returns an error with
// relevant information.
func (con *Connector) ExecQueryRow(tx *sql.Tx, query string, args ...interface{}) (*sql.Row, error) {
	var row *sql.Row

	// Prepare the SQL query.
	stm, err := tx.Prepare(query)
	if err != nil {
		return row, fmt.Errorf("can't prepare query: %w", err)
	}

	// Create a *sql.Row with the prepared statement and arguments.
	row = stm.QueryRow(args...)

	return row, nil
}

// FactoryConnector creates and returns a new *Connector by opening a connection
// to the SQL database using the provided Conf configuration. It returns the
// initialized Connector and an error if the connection cannot be established.
// The function uses the database driver and Data Source Name (DSN) specified
// in the configuration.
func FactoryConnector(c Conf) (*Connector, error) {
	var err error
	con := new(Connector)

	// Open a connection to the SQL database.
	con.DB, err = sql.Open(c.Driver, c.DSN)
	if err != nil {
		return con, fmt.Errorf("can't open connection to SqlDB(driver: %s): %w", c.Driver, err)
	}

	return con, nil
}
