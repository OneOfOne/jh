package jh // import "go.oneofone.dev/jh"

import (
	"encoding/json"
	"strconv"
)

type ValueType uint8

const (
	Null ValueType = iota
	Bool
	Number
	String
	Array
	Object
)

type Value struct {
	v json.RawMessage
}

func (v *Value) Type() ValueType {
	if len(v.v) == 0 {
		return Null
	}
	switch v.v[0] {
	case 't', 'T', 'f', 'F':
		return Bool
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return Number
	case '"':
		return String
	case '[':
		return Array
	case '{':
		return Object
	case 'n', 'N':
		fallthrough
	default:
		return Null
	}
}

func (v *Value) String() string {
	vv := v.v
	if len(vv) == 0 {
		return ""
	}
	if vv[0] == '"' && vv[len(vv)-1] == '"' {
		return string(vv[1 : len(vv)-1])
	}
	return string(vv)
}

func (v *Value) Int(base int) int64 {
	if base == 0 {
		base = 10
	}
	n, _ := strconv.ParseInt(v.String(), 10, 64)
	return n
}

func (v *Value) Uint(base int) uint64 {
	if base == 0 {
		base = 10
	}
	n, _ := strconv.ParseUint(v.String(), 10, 64)
	return n
}

func (v *Value) Float() float64 {
	n, _ := strconv.ParseFloat(v.String(), 64)
	return n
}

func (v *Value) Array() []*Value {
	var out []*Value
	v.As(&out)
	return out
}

func (v *Value) Object() map[string]*Value {
	var out map[string]*Value
	v.As(&out)
	return out
}

func (v *Value) As(ptr interface{}) error {
	return json.Unmarshal(v.v, ptr)
}
