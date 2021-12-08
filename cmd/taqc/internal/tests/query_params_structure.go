package tests

import "time"

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

//go:generate sh -c "go run $(cd ./\"$(git rev-parse --show-cdup)\" || exit; pwd)/cmd/taqc/taqc.go --type=TimeQueryParametersStructure"
type TimeQueryParametersStructure struct {
	Time      time.Time   `taqc:"time"`
	TimePtr   *time.Time  `taqc:"timePtr"`
	TimeSlice []time.Time `taqc:"timeSlice"`
}

//go:generate sh -c "go run $(cd ./\"$(git rev-parse --show-cdup)\" || exit; pwd)/cmd/taqc/taqc.go --type=PrimitiveTimeQueryParametersStructure"
type PrimitiveTimeQueryParametersStructure struct {
	UnixSec            time.Time `taqc:"sec"`
	UnixSec2           time.Time `taqc:"sec2, unixTimeUnit=sec"`
	UnixMilliSec       time.Time `taqc:"millisec, unixTimeUnit=millisec"`
	UnixMicroSec       time.Time `taqc:"microsec, unixTimeUnit=microsec"`
	UnixNanoSec        time.Time `taqc:"nanosec, unixTimeUnit=nanosec"`
	RFC3339            time.Time `taqc:"rfc3339, timeLayout=2006-01-02T15:04:05Z07:00"`
	PrioritizedRFC3339 time.Time `taqc:"prioritized_rfc3339, unixTimeUnit=sec, timeLayout=2006-01-02T15:04:05Z07:00"`
	//                                                     must be prioritized ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
}

//go:generate sh -c "go run $(cd ./\"$(git rev-parse --show-cdup)\" || exit; pwd)/cmd/taqc/taqc.go --type=PointerTimeQueryParametersStructure"
type PointerTimeQueryParametersStructure struct {
	UnixSec            *time.Time `taqc:"sec"`
	UnixSec2           *time.Time `taqc:"sec2, unixTimeUnit=sec"`
	UnixMilliSec       *time.Time `taqc:"millisec, unixTimeUnit=millisec"`
	UnixMicroSec       *time.Time `taqc:"microsec, unixTimeUnit=microsec"`
	UnixNanoSec        *time.Time `taqc:"nanosec, unixTimeUnit=nanosec"`
	RFC3339            *time.Time `taqc:"rfc3339, timeLayout=2006-01-02T15:04:05Z07:00"`
	PrioritizedRFC3339 *time.Time `taqc:"prioritized_rfc3339, unixTimeUnit=sec, timeLayout=2006-01-02T15:04:05Z07:00"`
	//                                                      must be prioritized ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
}

//go:generate sh -c "go run $(cd ./\"$(git rev-parse --show-cdup)\" || exit; pwd)/cmd/taqc/taqc.go --type=SliceTimeQueryParametersStructure"
type SliceTimeQueryParametersStructure struct {
	UnixSec            []time.Time `taqc:"sec"`
	UnixSec2           []time.Time `taqc:"sec2, unixTimeUnit=sec"`
	UnixMilliSec       []time.Time `taqc:"millisec, unixTimeUnit=millisec"`
	UnixMicroSec       []time.Time `taqc:"microsec, unixTimeUnit=microsec"`
	UnixNanoSec        []time.Time `taqc:"nanosec, unixTimeUnit=nanosec"`
	RFC3339            []time.Time `taqc:"rfc3339, timeLayout=2006-01-02T15:04:05Z07:00"`
	PrioritizedRFC3339 []time.Time `taqc:"prioritized_rfc3339, unixTimeUnit=sec, timeLayout=2006-01-02T15:04:05Z07:00"`
	//                                                       must be prioritized ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
}
