package providers

import (
	"context"
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type ActionFunc func(ctx *context.Context, input interface{}) interface{}
type ActionProvider interface {
	ActionsFuncs() map[string]ActionFunc
	Init()
}

func JsonParametersToProviderInputType[T any](data string) (*T, error) {
	var result T
	var rmap interface{}
	err := json.Unmarshal([]byte(data), &rmap)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(rmap, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ProviderResultToMap(result interface{}) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	err := mapstructure.Decode(result, &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap, nil
}
