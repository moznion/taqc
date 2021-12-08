package tests

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimitiveQueryParamsStructure_ToQueryParameters(t *testing.T) {
	q := &PrimitiveQueryParamsStructure{
		Foo:             "str",
		Bar:             123,
		Buz:             456.789,
		Qux:             true,
		ShouldBeIgnored: "should-be-ignored",
	}
	qp := q.ToQueryParameters()
	assert.EqualValues(t, url.Values{
		"foo": []string{"str"},
		"bar": []string{"123"},
		"buz": []string{"456.789000"},
		"qux": []string{"1"},
	}, qp)

	q = &PrimitiveQueryParamsStructure{
		Foo:             "str",
		Bar:             123,
		Buz:             456.789,
		Qux:             false,
		ShouldBeIgnored: "should-be-ignored",
	}
	qp = q.ToQueryParameters()
	assert.EqualValues(t, url.Values{
		"foo": []string{"str"},
		"bar": []string{"123"},
		"buz": []string{"456.789000"},
	}, qp)
}

func TestPointerQueryParamsStructure_ToQueryParameters(t *testing.T) {
	q := &PointerQueryParamsStructure{
		Foo: func() *string {
			v := "str"
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
		ShouldBeIgnored: func() *string {
			v := "should-be-ignored"
			return &v
		}(),
	}
	qp := q.ToQueryParameters()
	assert.EqualValues(t, url.Values{
		"foo": []string{"str"},
		"bar": []string{"123"},
		"buz": []string{"456.789000"},
		"qux": []string{"1"},
	}, qp)

	q = &PointerQueryParamsStructure{
		Foo: func() *string {
			v := "str"
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
			v := false
			return &v
		}(),
		ShouldBeIgnored: func() *string {
			v := "should-be-ignored"
			return &v
		}(),
	}
	qp = q.ToQueryParameters()
	assert.EqualValues(t, url.Values{
		"foo": []string{"str"},
		"bar": []string{"123"},
		"buz": []string{"456.789000"},
	}, qp)

	q = &PointerQueryParamsStructure{
		Foo:             nil,
		Bar:             nil,
		Buz:             nil,
		Qux:             nil,
		ShouldBeIgnored: nil,
	}
	qp = q.ToQueryParameters()
	assert.EqualValues(t, url.Values{}, qp)
}

func TestSliceQueryParamsStructure_ToQueryParameters(t *testing.T) {
	q := &SliceQueryParamsStructure{
		Foo:             []string{"str", "value"},
		Bar:             []int64{123, 456},
		Buz:             []float64{123.456, 234.567},
		ShouldBeIgnored: []string{"should", "be", "ignored"},
	}
	qp := q.ToQueryParameters()
	assert.EqualValues(t, url.Values{
		"foo": []string{"str", "value"},
		"bar": []string{"123", "456"},
		"buz": []string{"123.456000", "234.567000"},
	}, qp)
}
