# taqc 🚕

<b>Ta</b>g based url <b>Q</b>uery parameters <b>C</b>onstructor.  
(This is pronounced as same as "taxi")

## Synopsis

```go
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

fmt.Printf(queryParams.Encode())

// Output:
// buz=123&foo=string_value&foobar=1&qux=123.456000&qux=234.567000
```

## Description

This library constructs the query parameters (i.e. `url.Value{}`) according to the struct, and the `taqc` tag which is in each field.

Currently, it supports the following field types: `string`, `int64`, `float64`, `bool`, `*string`, `*int64`, `*float64`, `*bool`, `[]string`, `[]int64`, and `[]float64`.

If the bool field is `true`, the query parameter becomes `param_name=1`. Else, it omits the parameter.

And, when the pointer value is `nil`, it omits the parameter.

[![GoDoc](https://godoc.org/github.com/moznion/taqc?status.svg)](https://godoc.org/github.com/moznion/taqc)

## Command-line Tool

This library also provides a command-line tool to generate code.

`taqc.ConvertToQueryParams()` converts the given struct value to query parameters by using reflection dynamically,
but this way has a little disadvantage from the performance perspective.

This CLI tool generates the code statically, and this can convert that structure to the query parameters more efficiently.

### Installation

```
go install github.com/moznion/taqc/cmd/taqc@latest
```

### Usage

```
Usage of taqc:
  -type string
        [mandatory] a type name
  -output string
        [optional] output file name (default "srcdir/<type>_gen.go")
  -version
        show the version information
```

### Example

When it has the following code:

```go
//go:generate taqc --type=QueryParam
type QueryParam struct {
	Foo             string  `taqc:"foo"`
	Bar             int64   `taqc:"bar"`
	Buz             float64 `taqc:"buz"`
	Qux             bool    `taqc:"qux"`
}
```

then you run `go generate ./...`, it generates code on `query_param_gen.go` that is in the same directory of the original struct file.

## Author

moznion (<moznion@mail.moznion.net>)
