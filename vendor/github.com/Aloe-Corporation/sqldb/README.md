# Sqldb

The sqldb module defines a supercharged native `sql.DB` driver that simplifies configuration and add some usefull methods. It aims to minify code repetition that the `database/sql"` module can create.

![tests](https://github.com/Aloe-Corporation/sqldb/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/Aloe-Corporation/sqldb.svg)](https://pkg.go.dev/github.com/Aloe-Corporation/sqldb)

## Overview

The sqldb module offers:
- Simplified configuration
- A TryConnection method for pinging your database, triggering an error if the timeout is exceeded
- Convenient methods to streamline SQL database queries

## Concepts

Easily configure your connector with the provided `Conf` structure and its factory methods:

```go
type Conf struct {
	Driver string `yaml:"driver"` // Database driver (e.g., "mysql", "postgres", etc.).
	DSN    string `yaml:"dsn"`    // Data Source Name (DSN) for connecting to the database.
}
```
## Usage

### Configuration

The `sqldb.Conf` uses a YAML tags, it's easy to load SqlDB config with configuration file in your project

```go
type Conf struct {
	Driver string `yaml:"driver"` // example: postgres, mysql
	DSN    string `yaml:"dsn"` // connection string (format depends on the driver, read the associated documentation)
}
```

Example DSN:

- `postgres:` user=postgres password=example dbname=postgres host=localhost port=5432 sslmode=disable TimeZone=UTC
- `mysql:` root:example@tcp(localhost:3306)/dbtest?loc=UTC&tls=false&parseTime=true `(WARNING: parseTime=true is require)`


### Create new connector

To create new SqlDB Connector use this function with as configuration the structure `sqldb.FactoryConnector(c sqldb.Conf) (*sqldb.Connector, error)` and try connection with `sqldb.Connector.TryConnection(t int) err`

```go
var config = sqldb.Conf{
	Driver:     "mysql",
	DSN:      	"root:example@tcp(localhost:13306)/dbtest?loc=UTC&tls=false&parseTime=trueword",
}

// Build Connector
connector, err = sqldb.FactoryConnector(config)
if err != nil {
	return fmt.Errorf("fail to init SqlDB connector: %w", err)
}

// Test connection
err = connector.TryConnection(10)
if err != nil {
	return fmt.Errorf("fail to ping SqlDB: %w", err)
}

```
## Contributing

This section will be added soon.

## License

Sqldb module is released under the MIT license. See [LICENSE.txt](./LICENSE).
