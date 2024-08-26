package csv

import (
	"reflect"
)

type Marshaler interface {
	MarshalCSV() (string, error)
}

type Unmarshaler interface {
	UnmarshalCSV(string) error
}

var unmarshaler = reflect.TypeFor[Unmarshaler]()

type Ordered interface {
	OrderedCSV() []string
}
