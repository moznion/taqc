package taqc

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

var (
	ErrNilValueGiven             = errors.New("given value is nil")
	ErrQueryParameterNameIsEmpty = errors.New("query parameter name is empty in a tag")
	ErrUnsupportedFieldType      = errors.New("unsupported filed type has come")
)

const tagName = "taqc"

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
		tagValue, ok := tag.Lookup(tagName)
		if !ok { // nothing to do
			continue
		}

		splitTagValue := strings.Split(tagValue, ",")
		paramName := strings.TrimSpace(splitTagValue[0])
		if paramName == "" {
			return nil, ErrQueryParameterNameIsEmpty
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
