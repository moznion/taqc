package taqc

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/moznion/taqc/internal"
)

var (
	ErrNilValueGiven             = errors.New("given value is nil")
	ErrQueryParameterNameIsEmpty = errors.New("query parameter name is empty in a tag")
	ErrUnsupportedFieldType      = errors.New("unsupported filed type has come")
	ErrUnsupportedUnixTimeUnit   = errors.New("unsupported unix time unit has given")
)

// ConvertToQueryParams converts given structure to the query parameters according to the custom tags.
//
// When a field of the structure has `taqc` tag, it converts a value of that field to query parameter.
// Currently, it supports the following field types: `string`, `int64`, `float64`, `bool`, `*string`, `*int64`, `*float64`, `*bool`, `[]string`, `[]int64`, `[]float64`, `time.Time`, `*time.Time`, and `[]time.Time`.
// If the bool field is `true`, the query parameter becomes `param_name=1`. Else, it omits the parameter.
// And when the pointer value is `nil`, it omits the parameter.
//
// This library supports the `time.Time` fields. By default, it encodes that timestamp by `Time#Unix()`.
// If you want to encode it by another unix time format, you can use `unixTimeUnit` custom tag value.
// For example:
//
// 	type Query struct {
// 		Foo time.Time `taqc:"foo, unixTimeUnit=millisec"`
// 	}
//
// in the case of the above example, it encodes the timestamp by `Time#UnixMilli()`.
//
// Currently `unixTimeUnit` supports the following values:
//
// - `sec`
// - `millisec`
// - `microsec`
// - `nanosec`
//
// It also supports encoding with arbitrary time layout by `timeLayout` custom tag value. e.g.
//
// 	type Query struct {
// 		Foo time.Time `taqc:"foo, timeLayout=2006-01-02T15:04:05Z07:00"` // RFC3339 layout
// 	}
//
// then, it encodes the timestamp by `Time#Format()` with given layout.
//
// NOTE: `timeLayout` takes priority over `unixTimeUnit`. This means it uses `timeLayout` option even if you put them together.
func ConvertToQueryParams(v interface{}) (url.Values, error) {
	if v == nil {
		return nil, ErrNilValueGiven
	}

	qp := url.Values{}

	rv := reflect.ValueOf(v)
	elem := rv.Elem()
	for i := 0; i < elem.NumField(); i++ {
		typeField := elem.Type().Field(i)
		tag := typeField.Tag
		tagValue, ok := tag.Lookup(internal.TagName)
		if !ok { // nothing to do
			continue
		}

		splitTagValues := strings.Split(tagValue, ",")
		paramName := strings.TrimSpace(splitTagValues[0])
		if paramName == "" {
			return nil, ErrQueryParameterNameIsEmpty
		}

		timeLayout, unixTimeUnit := internal.ExtractTimeTag(splitTagValues[1:])

		timeFormatter := func(t time.Time) string {
			return fmt.Sprintf("%d", t.Unix())
		}
		if unixTimeUnit != "" {
			unixTimeGetter, err := getUnixTimeGetter(unixTimeUnit)
			if err != nil {
				return nil, err
			}
			timeFormatter = func(t time.Time) string {
				return fmt.Sprintf("%d", unixTimeGetter(t))
			}
		}
		if timeLayout != "" { // higher priority
			timeFormatter = func(t time.Time) string {
				return t.Format(timeLayout)
			}
		}

		field := elem.Field(i)
		fieldKind := field.Kind()
		switch fieldKind {
		case reflect.String:
			qp.Set(paramName, field.String())
		case reflect.Int64:
			qp.Set(paramName, fmt.Sprintf("%d", field.Int()))
		case reflect.Float64:
			qp.Set(paramName, fmt.Sprintf("%f", field.Float()))
		case reflect.Bool:
			if field.Bool() {
				qp.Set(paramName, "1")
			}
		case reflect.Struct:
			if field.Type().PkgPath() == "time" && field.Type().Name() == "Time" {
				t := field.Interface().(time.Time)
				qp.Set(paramName, timeFormatter(t))
			} else {
				return nil, fmt.Errorf("field type is %s: %w", fieldKind, ErrUnsupportedFieldType)
			}
		case reflect.Ptr:
			if !field.IsNil() {
				fieldElem := field.Elem()
				switch fieldElem.Kind() {
				case reflect.String:
					qp.Set(paramName, fieldElem.String())
				case reflect.Int64:
					qp.Set(paramName, fmt.Sprintf("%d", fieldElem.Int()))
				case reflect.Float64:
					qp.Set(paramName, fmt.Sprintf("%f", fieldElem.Float()))
				case reflect.Bool:
					if fieldElem.Bool() {
						qp.Set(paramName, "1")
					}
				case reflect.Struct:
					if fieldElem.Type().PkgPath() == "time" && fieldElem.Type().Name() == "Time" {
						t := fieldElem.Interface().(time.Time)
						qp.Set(paramName, timeFormatter(t))
					} else {
						return nil, fmt.Errorf("field type is *%s: %w", fieldKind, ErrUnsupportedFieldType)
					}
				default:
					return nil, fmt.Errorf("field type is *%s: %w", fieldKind, ErrUnsupportedFieldType)
				}
			}
		case reflect.Slice:
			l := field.Len()
			for j := 0; j < l; j++ {
				item := field.Index(j)
				switch item.Kind() {
				case reflect.String:
					qp.Add(paramName, item.String())
				case reflect.Int64:
					qp.Add(paramName, fmt.Sprintf("%d", item.Int()))
				case reflect.Float64:
					qp.Add(paramName, fmt.Sprintf("%f", item.Float()))
				case reflect.Struct:
					if item.Type().PkgPath() == "time" && item.Type().Name() == "Time" {
						t := item.Interface().(time.Time)
						qp.Add(paramName, timeFormatter(t))
					} else {
						return nil, fmt.Errorf("field type is []%s: %w", fieldKind, ErrUnsupportedFieldType)
					}
				default:
					return nil, fmt.Errorf("field type is []%s: %w", fieldKind, ErrUnsupportedFieldType)
				}
			}
		default:
			return nil, fmt.Errorf("field type is %s: %w", fieldKind, ErrUnsupportedFieldType)
		}
	}

	return qp, nil
}

func getUnixTimeGetter(unixTimeUnit string) (func(t time.Time) int64, error) {
	switch unixTimeUnit {
	case "", "sec":
		return func(t time.Time) int64 {
			return t.Unix()
		}, nil
	case "millisec":
		return func(t time.Time) int64 {
			return t.UnixMilli()
		}, nil
	case "microsec":
		return func(t time.Time) int64 {
			return t.UnixMicro()
		}, nil
	case "nanosec":
		return func(t time.Time) int64 {
			return t.UnixNano()
		}, nil
	default:
		return nil, fmt.Errorf("%s is unsupported: %w", unixTimeUnit, ErrUnsupportedUnixTimeUnit)
	}
}
