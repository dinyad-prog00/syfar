package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	t "syfar/types"

	"github.com/mitchellh/mapstructure"
)

type ActionFunc func(ctx *context.Context, input interface{}) (interface{}, error)

type Action struct {
	ActionFunc ActionFunc
	Inputs     []t.Input
}
type ActionProvider interface {
	GetActions() map[string]Action
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

func ValidateInput(name string, value interface{}, inputs []t.Input, pos string) error {
	for _, input := range inputs {
		if input.Name == name {
			if v := reflect.ValueOf(value); v.Kind() != input.Type {
				return fmt.Errorf("error at %s; type not mach, expect \"%s\" but got \"%s\" for \"%s\"", pos, input.Type.String(), v.Kind().String(), name)
			} else {
				return nil
			}
		}
	}

	return fmt.Errorf("error at %s; argument \"%s\" is not expected here", pos, name)
}
