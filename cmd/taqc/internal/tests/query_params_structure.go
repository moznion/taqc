package tests

//go:generate sh -c "go run $(cd ./\"$(git rev-parse --show-cdup)\" || exit; pwd)/cmd/taqc/taqc.go --type=PrimitiveQueryParamsStructure"
type PrimitiveQueryParamsStructure struct {
	Foo             string  `taqc:"foo"`
	Bar             int64   `taqc:"bar"`
	Buz             float64 `taqc:"buz"`
	Qux             bool    `taqc:"qux"`
	ShouldBeIgnored string
}

//go:generate sh -c "go run $(cd ./\"$(git rev-parse --show-cdup)\" || exit; pwd)/cmd/taqc/taqc.go --type=PointerQueryParamsStructure"
type PointerQueryParamsStructure struct {
	Foo             *string  `taqc:"foo"`
	Bar             *int64   `taqc:"bar"`
	Buz             *float64 `taqc:"buz"`
	Qux             *bool    `taqc:"qux"`
	ShouldBeIgnored *string
}

//go:generate sh -c "go run $(cd ./\"$(git rev-parse --show-cdup)\" || exit; pwd)/cmd/taqc/taqc.go --type=SliceQueryParamsStructure"
type SliceQueryParamsStructure struct {
	Foo             []string  `taqc:"foo"`
	Bar             []int64   `taqc:"bar"`
	Buz             []float64 `taqc:"buz"`
	ShouldBeIgnored []string
}
