package ginparamsmapper

import "fmt"

type MissingParamError struct {
	Param    string
	Location string
}

func (e MissingParamError) Error() string {
	return fmt.Sprintf("missing %v param %v in the context", e.Location, e.Param)
}

type NotCastableParamError struct {
	Param string
	Type  string
	Value string
	Err   error
}

func (e NotCastableParamError) Error() string {
	return fmt.Sprintf("%v (%v) not castable into %v. Error : "+e.Err.Error(), e.Param, e.Value, e.Type)
}
