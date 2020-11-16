package jh // import "go.oneofone.dev/jh"

import (
	"encoding/json"
	"strconv"
	"time"
)

type (
	Array  = []Value
	Object = map[string]Value
)

const (
	minMS = int64(1e12)
	minNS = int64(1e18)
)

var DefaultDateTimeLayouts = [...]string{time.RFC3339Nano, time.RFC1123, time.RFC1123Z, "2006-01-02 15:04:05", "2006-01-02", "2006/01/02"}

type Kind uint8

func (k Kind) String() string {
	switch k {
	case NullKind:
		return "null"
	case BoolKind:
		return "boolean"
	case NumberKind:
		return "number"
	case StringKind:
		return "string"
	case ArrayKind:
		return "array"
	case ObjectKind:
		return "object"
	default:
		return "unknown: " + strconv.Itoa(int(k))
	}
}

const (
	NullKind Kind = iota
	BoolKind
	NumberKind
	StringKind
	ArrayKind
	ObjectKind
)

type Value struct {
	v json.RawMessage
}

func (v Value) Kind() Kind {
	if len(v.v) == 0 {
		return NullKind
	}
	b := v.v[0]
	if b >= '0' && b <= '9' {
		return NumberKind
	}

	switch b {
	case 't', 'T', 'f', 'F':
		return BoolKind
	case '"':
		return StringKind
	case '[':
		return ArrayKind
	case '{':
		return ObjectKind
	case 'n', 'N':
		fallthrough
	default:
		return NullKind
	}
}

func isn(b byte) bool {
	return b >= '0' && b <= '9'
}

func (v Value) IsNull() bool {
	return v.Kind() == NullKind
}

func (v Value) String() string {
	vv := v.v
	if len(vv) == 0 {
		return ""
	}
	if vv[0] == '"' && vv[len(vv)-1] == '"' {
		return string(vv[1 : len(vv)-1])
	}
	return string(vv)
}

func (v Value) Bool() bool {
	b, _ := strconv.ParseBool(v.String())
	return b
}

func (v Value) Int(base int) int64 {
	if base == 0 {
		base = 10
	}
	n, _ := strconv.ParseInt(v.String(), 10, 64)
	return n
}

func (v Value) Uint(base int) uint64 {
	if base == 0 {
		base = 10
	}
	n, _ := strconv.ParseUint(v.String(), 10, 64)
	return n
}

func (v Value) Float() float64 {
	n, _ := strconv.ParseFloat(v.String(), 64)
	return n
}

func (v Value) Array() []*Value {
	var out []*Value
	v.As(&out)
	return out
}

func (v Value) Object() map[string]*Value {
	var out map[string]*Value
	v.As(&out)
	return out
}

// AsTime will try to return the time representation of the value, using the given fmts or DefaultDateTimeLayouts.
// if the value is a number it'll check if it's in NS, MS or a normal *U*nix timestamp and return that.
func (v Value) AsTime(fmts ...string) (fmt string, t time.Time, err error) {
	if len(fmts) == 0 {
		fmts = DefaultDateTimeLayouts[:]
	}
	if v.Kind() == NumberKind {
		n := v.Int(10)
		switch {
		case n >= minNS:
			return "NS", time.Unix(0, n), nil
		case n >= minNS:
			return "MS", time.Unix(n/1000, 0), nil
		default:
			return "U", time.Unix(n, 0), nil
		}
	}

	sv := v.String()
	for _, fmt = range fmts {
		if t, err = time.Parse(fmt, sv); err == nil {
			return
		}
	}

	return
}

func (v Value) As(ptr interface{}) error {
	return json.Unmarshal(v.v, ptr)
}

func (v Value) MarshalJSON() ([]byte, error) {
	return v.v.MarshalJSON()
}

func (v *Value) UnmarshalJSON(p []byte) error {
	return v.v.UnmarshalJSON(p)
}
