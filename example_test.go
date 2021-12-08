package taqc

import "fmt"

func ExampleConvertToQueryParams() {
	type Query struct {
		Foo             string    `taqc:"foo"`
		Bar             *string   `taqc:"bar"`
		Buz             int64     `taqc:"buz"`
		Qux             []float64 `taqc:"qux"`
		FooBar          bool      `taqc:"foobar"`
		Falsy           bool      `taqc:"falsy"`
		ShouldBeIgnored string
	}

	queryParams, err := ConvertToQueryParams(&Query{
		Foo:             "string_value",
		Bar:             nil, // <= should be ignored
		Buz:             123,
		Qux:             []float64{123.456, 234.567},
		FooBar:          true,  // <= be "foobar=1"
		Falsy:           false, // <= should be ignored
		ShouldBeIgnored: "should be ignored",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", queryParams.Encode())

	// Output:
	// buz=123&foo=string_value&foobar=1&qux=123.456000&qux=234.567000
}
