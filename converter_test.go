package taqc

import (
	"fmt"
	"net/url"
	"regexp"
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
		Foo regexp.Regexp `taqc:"foo"`
	}

	_, err := ConvertToQueryParams(&Query{
		Foo: *regexp.MustCompile(""),
	})
	assert.ErrorIs(t, err, ErrUnsupportedFieldType)
}

func TestConvertToQueryParams_ShouldRaiseErrorWhenInvalidPointerValue(t *testing.T) {
	type Query struct {
		Foo *regexp.Regexp `taqc:"foo"`
	}

	_, err := ConvertToQueryParams(&Query{
		Foo: regexp.MustCompile(""),
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

func TestConvertToQueryParams_ShouldRaiseErrorWhenInvalidStructSliceValue(t *testing.T) {
	type Query struct {
		Foo []regexp.Regexp `taqc:"foo"`
	}

	_, err := ConvertToQueryParams(&Query{
		Foo: []regexp.Regexp{func() regexp.Regexp {
			r := regexp.MustCompile("")
			return *r
		}()},
	})
	assert.ErrorIs(t, err, ErrUnsupportedFieldType)
}

func TestConvertToQueryParams_WithTime(t *testing.T) {
	type Query struct {
		UnixSec            time.Time `taqc:"sec"`
		UnixSec2           time.Time `taqc:"sec2, unixTimeUnit=sec"`
		UnixMilliSec       time.Time `taqc:"millisec, unixTimeUnit=millisec"`
		UnixMicroSec       time.Time `taqc:"microsec, unixTimeUnit=microsec"`
		UnixNanoSec        time.Time `taqc:"nanosec, unixTimeUnit=nanosec"`
		RFC3339            time.Time `taqc:"rfc3339, timeLayout=2006-01-02T15:04:05Z07:00"`
		PrioritizedRFC3339 time.Time `taqc:"prioritized_rfc3339, unixTimeUnit=sec, timeLayout=2006-01-02T15:04:05Z07:00"`
		//                                                     must be prioritized ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	}

	now := time.Now()
	qp, err := ConvertToQueryParams(&Query{
		UnixSec:            now,
		UnixSec2:           now,
		UnixMilliSec:       now,
		UnixMicroSec:       now,
		UnixNanoSec:        now,
		RFC3339:            now,
		PrioritizedRFC3339: now,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, url.Values{
		"sec":                 []string{fmt.Sprintf("%d", now.Unix())},
		"sec2":                []string{fmt.Sprintf("%d", now.Unix())},
		"millisec":            []string{fmt.Sprintf("%d", now.UnixMilli())},
		"microsec":            []string{fmt.Sprintf("%d", now.UnixMicro())},
		"nanosec":             []string{fmt.Sprintf("%d", now.UnixNano())},
		"rfc3339":             []string{now.Format(time.RFC3339)},
		"prioritized_rfc3339": []string{now.Format(time.RFC3339)},
	}, qp)
}

func TestConvertToQueryParams_WithPointerTime(t *testing.T) {
	type Query struct {
		UnixSec            *time.Time `taqc:"sec"`
		UnixSec2           *time.Time `taqc:"sec2, unixTimeUnit=sec"`
		UnixMilliSec       *time.Time `taqc:"millisec, unixTimeUnit=millisec"`
		UnixMicroSec       *time.Time `taqc:"microsec, unixTimeUnit=microsec"`
		UnixNanoSec        *time.Time `taqc:"nanosec, unixTimeUnit=nanosec"`
		RFC3339            *time.Time `taqc:"rfc3339, timeLayout=2006-01-02T15:04:05Z07:00"`
		PrioritizedRFC3339 *time.Time `taqc:"prioritized_rfc3339, unixTimeUnit=sec, timeLayout=2006-01-02T15:04:05Z07:00"`
		//                                                      must be prioritized ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	}

	now := time.Now()
	qp, err := ConvertToQueryParams(&Query{
		UnixSec:            &now,
		UnixSec2:           &now,
		UnixMilliSec:       &now,
		UnixMicroSec:       &now,
		UnixNanoSec:        &now,
		RFC3339:            &now,
		PrioritizedRFC3339: &now,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, url.Values{
		"sec":                 []string{fmt.Sprintf("%d", now.Unix())},
		"sec2":                []string{fmt.Sprintf("%d", now.Unix())},
		"millisec":            []string{fmt.Sprintf("%d", now.UnixMilli())},
		"microsec":            []string{fmt.Sprintf("%d", now.UnixMicro())},
		"nanosec":             []string{fmt.Sprintf("%d", now.UnixNano())},
		"rfc3339":             []string{now.Format(time.RFC3339)},
		"prioritized_rfc3339": []string{now.Format(time.RFC3339)},
	}, qp)

	qp, err = ConvertToQueryParams(&Query{
		UnixSec:            nil,
		UnixSec2:           nil,
		UnixMilliSec:       nil,
		UnixMicroSec:       nil,
		UnixNanoSec:        nil,
		RFC3339:            nil,
		PrioritizedRFC3339: nil,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, url.Values{}, qp)
}

func TestConvertToQueryParams_WithTimeSlice(t *testing.T) {
	type Query struct {
		UnixSec            []time.Time `taqc:"sec"`
		UnixSec2           []time.Time `taqc:"sec2, unixTimeUnit=sec"`
		UnixMilliSec       []time.Time `taqc:"millisec, unixTimeUnit=millisec"`
		UnixMicroSec       []time.Time `taqc:"microsec, unixTimeUnit=microsec"`
		UnixNanoSec        []time.Time `taqc:"nanosec, unixTimeUnit=nanosec"`
		RFC3339            []time.Time `taqc:"rfc3339, timeLayout=2006-01-02T15:04:05Z07:00"`
		PrioritizedRFC3339 []time.Time `taqc:"prioritized_rfc3339, unixTimeUnit=sec, timeLayout=2006-01-02T15:04:05Z07:00"`
		//                                                       must be prioritized ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	}

	now := time.Now()
	qp, err := ConvertToQueryParams(&Query{
		UnixSec:            []time.Time{now, now},
		UnixSec2:           []time.Time{now, now},
		UnixMilliSec:       []time.Time{now, now},
		UnixMicroSec:       []time.Time{now, now},
		UnixNanoSec:        []time.Time{now, now},
		RFC3339:            []time.Time{now, now},
		PrioritizedRFC3339: []time.Time{now, now},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, url.Values{
		"sec":                 []string{fmt.Sprintf("%d", now.Unix()), fmt.Sprintf("%d", now.Unix())},
		"sec2":                []string{fmt.Sprintf("%d", now.Unix()), fmt.Sprintf("%d", now.Unix())},
		"millisec":            []string{fmt.Sprintf("%d", now.UnixMilli()), fmt.Sprintf("%d", now.UnixMilli())},
		"microsec":            []string{fmt.Sprintf("%d", now.UnixMicro()), fmt.Sprintf("%d", now.UnixMicro())},
		"nanosec":             []string{fmt.Sprintf("%d", now.UnixNano()), fmt.Sprintf("%d", now.UnixNano())},
		"rfc3339":             []string{now.Format(time.RFC3339), now.Format(time.RFC3339)},
		"prioritized_rfc3339": []string{now.Format(time.RFC3339), now.Format(time.RFC3339)},
	}, qp)
}

func TestConvertToQueryParams_ShouldRaiseErrorWithInvalidUnixTimeUnit(t *testing.T) {
	type Query1 struct {
		InvalidUnixSec time.Time `taqc:"invalidUnixSec, unixTimeUnit=INVALID"`
	}
	_, err := ConvertToQueryParams(&Query1{
		InvalidUnixSec: time.Now(),
	})
	assert.ErrorIs(t, err, ErrUnsupportedUnixTimeUnit)

	type Query2 struct {
		InvalidUnixSec *time.Time `taqc:"invalidUnixSec, unixTimeUnit=INVALID"`
	}
	_, err = ConvertToQueryParams(&Query2{
		InvalidUnixSec: func() *time.Time {
			now := time.Now()
			return &now
		}(),
	})
	assert.ErrorIs(t, err, ErrUnsupportedUnixTimeUnit)

	type Query3 struct {
		InvalidUnixSec []time.Time `taqc:"invalidUnixSec, unixTimeUnit=INVALID"`
	}
	_, err = ConvertToQueryParams(&Query3{
		InvalidUnixSec: []time.Time{time.Now()},
	})
	assert.ErrorIs(t, err, ErrUnsupportedUnixTimeUnit)
}

func TestConvertToQueryParams_WithUnsupportedFieldType(t *testing.T) {
	type Query struct {
		Invalid interface{} `taqc:"invalid"`
	}
	_, err := ConvertToQueryParams(&Query{
		Invalid: struct{}{},
	})
	assert.ErrorIs(t, err, ErrUnsupportedFieldType)
}

func TestConvertToQueryParams_WithUnsupportedPointerFieldType(t *testing.T) {
	type Query struct {
		Invalid *uintptr `taqc:"invalid"`
	}
	ptr := uintptr(0)
	_, err := ConvertToQueryParams(&Query{
		Invalid: &ptr,
	})
	assert.ErrorIs(t, err, ErrUnsupportedFieldType)
}
