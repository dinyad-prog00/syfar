package types

import "reflect"

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

type TestResult struct {
	Id           int
	Passed       bool
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
