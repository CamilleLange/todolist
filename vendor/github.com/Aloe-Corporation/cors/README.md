# Cors
The Cors module enhances `github.com/gin-contrib/cors` with more convenient methods and configuration.

This project has been developped by the [Aloe](https://www.aloe-corp.com/) team and is now open source.

![tests](https://github.com/Aloe-Corporation/cors/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/Aloe-Corporation/cors.svg)](https://pkg.go.dev/github.com/Aloe-Corporation/cors)

## Overview

The cors module offers:
- A `Conf` structure with `yaml` and `mapstructure` tags
- Builder methods for the convenient creation of a new middleware
- Default values to reduce code repetition in your projects
- Easy integration because it embed a native `github.com/gin-contrib/cors`

## Concepts

### Build your middleware on default configuration

The default middleware provides an open CORS configuration that can be useful during the development process.

```go
    builder := new(cors.Builder)
    conf := builder.New().Build()
```

### Easily add custom values

Enhance the security of your Gin API by specifying allowed origins in addition to the default configuration.

```go
    builder := new(cors.Builder)
    conf := builder.New().WithOrigins("http://localhost:8080").Build()
```

## Usage

Use the cors middleware on your gin endpoints.

```go
    Router := gin.New()
    corsBuilder := new(cors.Builder)

    Router.Use(cors.Middleware(corsBuilder.New().WithOrigins("http://localhost:8080").Build()))
```

Use it with default configuration.

```go
    Router := gin.New()
    Router.Use(cors.Middleware(nil))
```

Use if with the `Conf` struct

```go
    Router := gin.New()

    conf := &cors.Conf{
        // your parsed configuration
    }

    corsBuilder := new(cors.Builder)
  
    Router.Use(cors.Middleware(corsBuilder.NewFromConf(conf).Build()))
```

## Contributing

This section will be added soon.

## License

Client is released under the MIT license. See [LICENSE.txt](./LICENSE).