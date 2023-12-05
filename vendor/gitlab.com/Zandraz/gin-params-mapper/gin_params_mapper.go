package ginparamsmapper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ParamLocation string

const (
	ContextVarLocation ParamLocation = "context"
	QueryLocation      ParamLocation = "query"
	PathLocation       ParamLocation = "path"
)

var (
	ErrUnknownParamLocation = errors.New("unknown param location")
)

// GetFromContext a param located at the given loc.
// The param will be converted using the dest type and convertion result will be set into dest.
func GetFromContext(loc ParamLocation, name string, c *gin.Context, dest any) error {
	switch loc {
	case ContextVarLocation:
		return GetVarFromContext(name, c, dest)
	case QueryLocation:
		return GetQueryParamFromContext(name, c, dest)
	case PathLocation:
		return GetPathParamFromContext(name, c, dest)
	default:
		return ErrUnknownParamLocation
	}
}

// GetQueryParamFromContext using the query param name, from the given *gin.Context.
// The param will be converted using the dest type and convertion result will be set into dest.
func GetQueryParamFromContext(name string, c *gin.Context, dest any) error {
	param, exists := c.GetQuery(name)
	if !exists {
		return MissingParamError{
			Param:    name,
			Location: string(QueryLocation),
		}
	}
	return convert(param, name, dest)
}

// GetPathParamFromContext using the path param name, from the given *gin.Context.
// The param will be converted using the dest type and convertion result will be set into dest.
func GetPathParamFromContext(name string, c *gin.Context, dest any) error {
	param := c.Param(name)
	if param == "" {
		return MissingParamError{
			Param:    name,
			Location: string(PathLocation),
		}
	}
	return convert(param, name, dest)
}

// GetVarFromContext using the variable name, from the given *gin.Context.
// The context variable will be converted using the dest type and convertion result will be set into dest.
func GetVarFromContext(name string, c *gin.Context, dest any) error {
	variable, exist := c.Get(name)
	if !exist {
		return MissingParamError{
			Param:    name,
			Location: string(ContextVarLocation),
		}
	}
	return convert(variable, name, dest)
}

func convert(param any, name string, dest any) error {
	s, err := interfaceToString(param)
	if err != nil {
		return fmt.Errorf("can't convert param to string. Error : %w", err)
	}

	switch d := dest.(type) {
	case *string:
		*d = s

	case *float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return NotCastableParamError{
				Param: name,
				Type:  "*float64",
				Value: s,
				Err:   err,
			}
		}
		*d = f

	case *int:
		i, err := strconv.Atoi(s)
		if err != nil {
			return NotCastableParamError{
				Param: name,
				Type:  "*int",
				Value: s,
				Err:   err,
			}
		}
		*d = i

	case *int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return NotCastableParamError{
				Param: name,
				Type:  "*int64",
				Value: s,
				Err:   err,
			}
		}
		*d = i

	case *bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return NotCastableParamError{
				Param: name,
				Type:  "*bool",
				Value: s,
				Err:   err,
			}
		}
		*d = b

	case *[]byte:
		*d = []byte(s)

	case *uuid.UUID:
		uid, err := uuid.Parse(s)
		if err != nil {
			return NotCastableParamError{
				Param: name,
				Type:  "uuid.UUID",
				Value: s,
				Err:   err,
			}
		}
		*d = uid

	case *time.Time:
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return NotCastableParamError{
				Param: name,
				Type:  "time.Time",
				Value: s,
				Err:   err,
			}
		}
		*d = t

	case *time.Duration:
		duration, err := time.ParseDuration(s)
		if err != nil {
			return NotCastableParamError{
				Param: name,
				Type:  "time.Duration",
				Value: s,
				Err:   err,
			}
		}
		*d = duration

	default:
		return NotCastableParamError{
			Param: name,
			Type:  fmt.Sprintf("%t", dest),
			Value: s,
			Err:   errors.New("destination type not supported"),
		}
	}
	return nil
}

func interfaceToString(i any) (s string, err error) {
	switch t := i.(type) {
	case string:
		s = t

	case float64:
		s = strconv.FormatFloat(t, 'f', 2, 64)

	case int:
		s = strconv.Itoa(t)

	case int64:
		s = strconv.FormatInt(t, 10)

	case bool:
		s = strconv.FormatBool(t)

	case []byte:
		s = string(t)

	case time.Time:
		s = t.Format(time.RFC3339Nano)

	// if you have to convert interfaces which implement fmt.Stringer in an other way,
	// be sure to put this specials cases before this one.
	case fmt.Stringer:
		s = t.String()

	default:
		err = fmt.Errorf("given interface has a unsupported type. Type : %t", t)
	}
	return
}
