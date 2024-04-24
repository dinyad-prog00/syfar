package types

import "reflect"

type TestState int8

const (
	StatePassed TestState = iota
	StateSkipped
	StateFailed
)

type ExpectationItemResult struct {
	Id      int
	Passed  bool
	Message string
}

type ExpectationResult struct {
	Id     int
	Passed bool
	Items  []ExpectationItemResult
}

type SyfarResult struct {
	TestsResult   []TestResult
	NbTestsPassed int
	NbTestsFailed int
	NbTestSkipped int
}

type TestResult struct {
	Id           int
	State        TestState
	Description  string
	Expectations []ExpectationResult
}

type ActionResult struct {
	Id          int
	Description string
}

type Input struct {
	Name     string
	Type     reflect.Kind
	Required bool
}

type ActionResultContext struct {
	Result        interface{}
	MapResult     map[string]interface{}
	DumpMapResult map[string]interface{}
}
