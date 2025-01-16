package csv

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"unsafe"

	stdcsv "encoding/csv"
)

type eface struct {
	typ unsafe.Pointer
	val unsafe.Pointer
}

func UnmarshalCSVReader[T any](reader *stdcsv.Reader) ([]T, error) {
	header, rerr := reader.Read()
	if rerr != nil {
		return nil, rerr
	}
	if len(header) != 0 {
		header[0] = strings.TrimPrefix(header[0], "\xef\xbb\xbf")
	}

	var i any = *new(T)
	ptr := uintptr((*eface)(unsafe.Pointer(&i)).typ)
	m := load(ptr)
	if m == nil {
		m = parse(reflect.TypeFor[T]())
		fieldsCache.Store(ptr, m)
	}

	fields := make([]*field, len(header))
	for i, h := range header {
		fields[i] = m.fields[h]
	}

	records := make([]T, 0, 16)

	for {
		items, rerr := reader.Read()
		if rerr == io.EOF {
			return records, nil
		}
		if rerr != nil {
			return nil, rerr
		}

		var record T
		elem := reflect.ValueOf(&record).Elem()

		for i, s := range items {
			field := fields[i]
			if field == nil {
				continue
			}
			f := elem.Field(field.index)
			if field.unmarshaler {
				if f.Kind() == reflect.Pointer && f.IsNil() {
					f.Set(reflect.New(f.Type().Elem()))
				} else if f.Kind() == reflect.Map && f.IsNil() {
					f.Set(reflect.MakeMap(f.Type()))
				}
				uerr := f.Interface().(Unmarshaler).UnmarshalCSV(s)
				if uerr != nil {
					return nil, uerr
				}
			} else if field.str {
				f.SetString(s)
			} else if field.bytes {
				f.SetBytes([]byte(s))
			}
		}

		records = append(records, record)
	}
}

func UnmarshalReader[T any](r io.Reader) ([]T, error) {
	reader := stdcsv.NewReader(r)
	reader.ReuseRecord = true
	return UnmarshalCSVReader[T](reader)
}

func UnmarshalFile[T any](filename string) (records []T, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	records, err = UnmarshalReader[T](f)
	ferr := f.Close()
	if ferr != nil {
		err = errors.Join(err, ferr)
	}
	return
}

func Unmarshal[T any](data []byte) (records []T, err error) {
	return UnmarshalReader[T](bytes.NewBuffer(data))
}
