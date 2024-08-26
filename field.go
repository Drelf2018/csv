package csv

import (
	"reflect"
	"sync"
)

type field struct {
	index       int
	str         bool
	bytes       bool
	unmarshaler bool
}

type fields struct {
	keys   []string
	fields map[string]*field
}

var fieldsCache sync.Map

func load(ptr uintptr) *fields {
	m, ok := fieldsCache.Load(ptr)
	if ok {
		return m.(*fields)
	}
	return nil
}

func parse(typ reflect.Type) *fields {
	num := typ.NumField()
	fields := &fields{make([]string, 0, num), make(map[string]*field, num)}
	for i := 0; i < num; i++ {
		f := typ.Field(i)
		tag := f.Tag.Get("csv")
		if tag != "" {
			fields.keys = append(fields.keys, tag)
			field := &field{
				index:       i,
				str:         f.Type.Kind() == reflect.String,
				bytes:       f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Uint8,
				unmarshaler: f.Type.Implements(unmarshaler),
			}
			if field.str || field.bytes || field.unmarshaler {
				fields.fields[tag] = field
			}
		}
	}
	return fields
}
