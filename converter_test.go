package taqc

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertToQueryParams(t *testing.T) {
	type Query struct {
		Foo             string  `taqc:"foo"`
		Bar             int64   `taqc:"bar"`
		Buz             float64 `taqc:"buz"`
		Qux             bool    `taqc:"qux"`
		FooBar          bool    `taqc:"foobar"`
		ShouldBeIgnored string
	}

	qp, err := ConvertToQueryParams(&Query{
		Foo:             "str-value",
		Bar:             123,
		Buz:             456.789,
		Qux:             true,
		FooBar:          false,
		ShouldBeIgnored: "should-be-ignored",
	})
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

	qp, err := ConvertToQueryParams(&Query{
		Foo: []string{"s1", "s2"},
		Bar: []int64{123, 456},
		Buz: []float64{123.456, 234.567},
	})
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

	qp, err := ConvertToQueryParams(&Query{
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
	})
	assert.NoError(t, err)
	assert.EqualValues(t, url.Values{
		"foo": []string{"str-value"},
		"bar": []string{"123"},
		"buz": []string{"456.789000"},
		"qux": []string{"1"},
	}, qp)

	qp, err = ConvertToQueryParams(&Query{
		Foo:    nil,
		Bar:    nil,
		Buz:    nil,
		Qux:    nil,
		FooBar: nil,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, url.Values{}, qp)
}

func TestConvertToQueryParams_ShouldRaiseErrorWhenNilValueGiven(t *testing.T) {
	_, err := ConvertToQueryParams(nil)
	assert.ErrorIs(t, err, ErrNilValueGiven)
}

func TestConvertToQueryParams_ShouldRaiseErrorWhenTagParamValueIsEmpty(t *testing.T) {
	type Query struct {
		Foo string `taqc:""`
	}

	_, err := ConvertToQueryParams(&Query{
		Foo: "foo",
	})
	assert.ErrorIs(t, err, ErrQueryParameterNameIsEmpty)
}

func TestConvertToQueryParams_ShouldRaiseErrorWhenInvalidPrimitiveValue(t *testing.T) {
	type Query struct {
		Foo time.Time `taqc:"foo"`
	}

	_, err := ConvertToQueryParams(&Query{
		Foo: time.Now(),
	})
	assert.ErrorIs(t, err, ErrUnsupportedFieldType)
}

func TestConvertToQueryParams_ShouldRaiseErrorWhenInvalidPointerValue(t *testing.T) {
	type Query struct {
		Foo *time.Time `taqc:"foo"`
	}

	_, err := ConvertToQueryParams(&Query{
		Foo: func() *time.Time {
			t := time.Now()
			return &t
		}(),
	})
	assert.ErrorIs(t, err, ErrUnsupportedFieldType)
}

func TestConvertToQueryParams_ShouldRaiseErrorWhenInvalidSliceValue(t *testing.T) {
	type Query struct {
		Foo []bool `taqc:"foo"`
	}

	_, err := ConvertToQueryParams(&Query{
		Foo: []bool{true},
	})
	assert.ErrorIs(t, err, ErrUnsupportedFieldType)
}
