package taqc

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToQueryParams(t *testing.T) {
	type Query struct {
		Foo    string  `taqc:"foo"`
		Bar    int64   `taqc:"bar"`
		Buz    float64 `taqc:"buz"`
		Qux    bool    `taqc:"qux"`
		FooBar bool    `taqc:"foobar"`
	}

	q := &Query{
		Foo:    "str-value",
		Bar:    123,
		Buz:    456.789,
		Qux:    true,
		FooBar: false,
	}
	qp, err := ConvertToQueryParams(q)
	assert.NoError(t, err)

	assert.EqualValues(t, url.Values{
		"foo": []string{"str-value"},
		"bar": []string{"123"},
		"buz": []string{"456.789000"},
		"qux": []string{"1"},
	}, qp)
}

func TestConvertToQueryParams_ForSliceFields(t *testing.T) {
	type Query struct {
		Foo []string  `taqc:"foo"`
		Bar []int64   `taqc:"bar"`
		Buz []float64 `taqc:"buz"`
	}

	q := &Query{
		Foo: []string{"s1", "s2"},
		Bar: []int64{123, 456},
		Buz: []float64{123.456, 234.567},
	}
	qp, err := ConvertToQueryParams(q)
	assert.NoError(t, err)

	assert.EqualValues(t, url.Values{
		"foo": []string{"s1", "s2"},
		"bar": []string{"123", "456"},
		"buz": []string{"123.456000", "234.567000"},
	}, qp)
}

func TestConvertToQueryParams_ForPointerFields(t *testing.T) {
	type Query struct {
		Foo    *string  `taqc:"foo"`
		Bar    *int64   `taqc:"bar"`
		Buz    *float64 `taqc:"buz"`
		Qux    *bool    `taqc:"qux"`
		FooBar *bool    `taqc:"foobar"`
	}

	q := &Query{
		Foo: func() *string {
			v := "str-value"
			return &v
		}(),
		Bar: func() *int64 {
			v := int64(123)
			return &v
		}(),
		Buz: func() *float64 {
			v := 456.789
			return &v
		}(),
		Qux: func() *bool {
			v := true
			return &v
		}(),
		FooBar: func() *bool {
			v := false
			return &v
		}(),
	}
	qp, err := ConvertToQueryParams(q)
	assert.NoError(t, err)

	assert.EqualValues(t, url.Values{
		"foo": []string{"str-value"},
		"bar": []string{"123"},
		"buz": []string{"456.789000"},
		"qux": []string{"1"},
	}, qp)
}
