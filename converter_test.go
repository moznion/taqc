package taqc

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToQueryParams(t *testing.T) {
	type Query struct {
		Foo string  `taqc:"foo"`
		Bar int64   `taqc:"bar"`
		Buz float64 `taqc:"buz"`
		Qux bool    `taqc:"qux"`
	}

	q := &Query{
		Foo: "str-value",
		Bar: 123,
		Buz: 456.789,
		Qux: true,
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
