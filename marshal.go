package csv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"unsafe"

	stdcsv "encoding/csv"
)

func MarshalCSVWriter[T any](writer *stdcsv.Writer, records []T) error {
	var i any = *new(T)
	ptr := uintptr((*eface)(unsafe.Pointer(&i)).typ)
	m := load(ptr)
	if m == nil {
		m = parse(reflect.TypeFor[T]())
		fieldsCache.Store(ptr, m)
	}

	var header []string
	if x, ok := i.(Ordered); ok {
		header = x.OrderedCSV()
	} else {
		header = m.keys
	}

	err := writer.Write(header)
	if err != nil {
		return err
	}

	fields := make([]*field, len(header))
	for i, h := range header {
		fields[i] = m.fields[h]
	}

	items := make([]string, len(header))

	for _, record := range records {
		val := reflect.ValueOf(record)
		for idx, field := range fields {
			if field == nil {
				continue
			}
			f := val.Field(field.index)
			switch i := f.Interface().(type) {
			case Marshaler:
				items[idx], err = i.MarshalCSV()
				if err != nil {
					return err
				}
			case string:
				items[idx] = i
			case []byte:
				items[idx] = string(i)
			default:
				items[idx] = fmt.Sprint(i)
			}
		}
		err = writer.Write(items)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}

func MarshalWriter[T any](r io.Writer, records []T) error {
	writer := stdcsv.NewWriter(r)
	return MarshalCSVWriter(writer, records)
}

func MarshalFile[T any](filename string, records []T) (err error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return
	}
	// 写入 UTF-8 BOM 避免使用  Microsoft Excel 打开乱码
	_, err = f.WriteString("\xef\xbb\xbf")
	if err != nil {
		return
	}
	err = MarshalWriter(f, records)
	ferr := f.Close()
	if ferr != nil {
		err = errors.Join(err, ferr)
	}
	return
}

func Marshal[T any](records []T) (p []byte, err error) {
	buf := &bytes.Buffer{}
	err = MarshalWriter(buf, records)
	p = buf.Bytes()
	return
}
