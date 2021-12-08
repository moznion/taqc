package tests

import (
	"fmt"
	"net/url"
	"testing"
	"time"

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

func TestTimeQueryParametersStructure_ToQueryParameters(t *testing.T) {
	now := time.Now()
	q := &TimeQueryParametersStructure{
		Time:      now,
		TimePtr:   &now,
		TimeSlice: []time.Time{now},
	}

	unixTimeStr := fmt.Sprintf("%d", now.Unix())

	qp := q.ToQueryParameters()
	assert.EqualValues(t, url.Values{
		"time":      []string{unixTimeStr},
		"timePtr":   []string{unixTimeStr},
		"timeSlice": []string{unixTimeStr},
	}, qp)
}

func TestPrimitiveTimeQueryParametersStructure_ToQueryParameters(t *testing.T) {
	now := time.Now()
	q := &PrimitiveTimeQueryParametersStructure{
		UnixSec:            now,
		UnixSec2:           now,
		UnixMilliSec:       now,
		UnixMicroSec:       now,
		UnixNanoSec:        now,
		RFC3339:            now,
		PrioritizedRFC3339: now,
	}

	qp := q.ToQueryParameters()
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

func TestPointerTimeQueryParametersStructure_ToQueryParameters(t *testing.T) {
	now := time.Now()
	q := &PointerTimeQueryParametersStructure{
		UnixSec:            &now,
		UnixSec2:           &now,
		UnixMilliSec:       &now,
		UnixMicroSec:       &now,
		UnixNanoSec:        &now,
		RFC3339:            &now,
		PrioritizedRFC3339: &now,
	}
	qp := q.ToQueryParameters()
	assert.EqualValues(t, url.Values{
		"sec":                 []string{fmt.Sprintf("%d", now.Unix())},
		"sec2":                []string{fmt.Sprintf("%d", now.Unix())},
		"millisec":            []string{fmt.Sprintf("%d", now.UnixMilli())},
		"microsec":            []string{fmt.Sprintf("%d", now.UnixMicro())},
		"nanosec":             []string{fmt.Sprintf("%d", now.UnixNano())},
		"rfc3339":             []string{now.Format(time.RFC3339)},
		"prioritized_rfc3339": []string{now.Format(time.RFC3339)},
	}, qp)

	q = &PointerTimeQueryParametersStructure{
		UnixSec:            nil,
		UnixSec2:           nil,
		UnixMilliSec:       nil,
		UnixMicroSec:       nil,
		UnixNanoSec:        nil,
		RFC3339:            nil,
		PrioritizedRFC3339: nil,
	}
	qp = q.ToQueryParameters()
	assert.EqualValues(t, url.Values{}, qp)
}

func TestSliceTimeQueryParametersStructure_ToQueryParameters(t *testing.T) {
	now := time.Now()
	q := &SliceTimeQueryParametersStructure{
		UnixSec:            []time.Time{now, now},
		UnixSec2:           []time.Time{now, now},
		UnixMilliSec:       []time.Time{now, now},
		UnixMicroSec:       []time.Time{now, now},
		UnixNanoSec:        []time.Time{now, now},
		RFC3339:            []time.Time{now, now},
		PrioritizedRFC3339: []time.Time{now, now},
	}
	qp := q.ToQueryParameters()
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
