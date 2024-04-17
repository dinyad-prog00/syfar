package assertions

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func deepEqual(a, b interface{}) bool {
	if !reflect.DeepEqual(a, b) {
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}
	return true
}

func areSameTypes(i, j interface{}) bool {
	if i == nil && j != nil || i != nil && j == nil {
		return false
	}

	var err error
	i, j, err = handleJSONNumber(i, j)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(
		reflect.Zero(reflect.TypeOf(i)).Interface(),
		reflect.Zero(reflect.TypeOf(j)).Interface(),
	)
}

func handleJSONNumber(actual interface{}, expected interface{}) (interface{}, interface{}, error) {
	jsNumber, is := actual.(json.Number)
	if !is {
		return actual, expected, nil
	}

	switch expected.(type) {
	case string:
		return jsNumber.String(), expected, nil
	case int64:
		i, err := jsNumber.Int64()
		if err != nil {
			return actual, expected, err
		}
		return i, expected, nil
	case float64:
		f, err := jsNumber.Float64()
		if err != nil {
			return actual, expected, err
		}
		return f, expected, nil
	}

	return jsNumber, expected, nil
}
