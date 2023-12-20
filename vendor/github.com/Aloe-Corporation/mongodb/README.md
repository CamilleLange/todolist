# mongodb

The mongodb module defines supercharged MongoDB Connector that adds some usefull functions to the native `go.mongodb.org/mongo-driver/mongo` client.

This project has been developped by the [Aloe](https://www.aloe-corp.com/) team and is now open source.

![tests](https://github.com/Aloe-Corporation/mongodb/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/Aloe-Corporation/mongodb.svg)](https://pkg.go.dev/github.com/Aloe-Corporation/mongodb)

## Overview

The mongodb connector supercharge the native `go.mongodb.org/mongo-driver/mongo` client.

The mongodb module provides:

- A structure to hold configuration
- A factory to create the connector
- Addionnal functions
- Infinite compatibility because it embed a native `go.mongodb.org/mongo-driver/mongo` client

## Concepts 

The module aims to ease the configuration and some operations to interact with a MongoDB server.

## Usage

### Configuration
The `mongodb.Conf` use a YAML tags, it's easy to load MongoDB config with configuration file in your project
```go
type Conf struct {
	DB         string `yaml:"db"`          // Name of the database.
	Host       string `yaml:"host"`        // URL to reach the mongoDB server.
	Port       int    `yaml:"port,omitempty"` // Optionnal port, if set to 0 it won't be processed.
	Username   string `yaml:"username"`    // Credential to authenticate to the db.
	Password   string `yaml:"password"`    // Credential to authenticate to the db.
	AuthSource string `yaml:"auth_source"` // Database to check authentication
	Timeout    int    `yaml:"timeout"`     // Connection timeout in seconds
}
```

### Create new connector
To create a new MongoDB Connector use this function with as configuration the structure `mongodb.FactoryConnector(c mongodb.Conf) (*mongodb.Connector, error)` and try connection with `mongodb.Connector.TryConnection() err`
```go
var config := mongodb.Conf{
	DB:       "my_database",
	Host:     "localhost:27006",
	Username: "user",
	Password: "pass",
	AuthSource: "admin"
	Timeout:  10,
}

md, err = mongodb.FactoryConnector(config)
if err != nil {
	return fmt.Errorf("fail to init MongoDB connector: %w", err)
}

err = md.TryConnection()
if err != nil {
	return fmt.Errorf("fail to ping MongoDB: %w", err)
}

```

## Test
To run test use:
- `make test`

All environment variables present in the `test/.env.example` must be set in your test environment.


## Contributing

This section will be added soon.

## License

Mongodb module is released under the MIT license. See [LICENSE.txt](./LICENSE).
