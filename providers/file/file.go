package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	pvd "syfar/providers"
	r "syfar/runner"
	t "syfar/types"
)

/**
Input config
*/

var ReadFileInput = []t.Input{
	{Name: "path", Type: reflect.String, Required: true},
}

type FileActionProvider struct {
}

func ReadFile(ctx *context.Context, params interface{}) (interface{}, error) {
	paramString, ok := params.(string)
	if !ok {
		return nil, fmt.Errorf("params arg should be a string")
	}

	input, err := pvd.JsonParametersToProviderInputType[FileProviderInput](paramString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rootDir := r.GetValueFromContext(*ctx, "syfar.rootdir").(string)
	filePath := filepath.Join(rootDir, input.Path)
	data, err := os.ReadFile(filePath)
	if err != nil {

		return nil, err

	}

	result := Result{}
	result.Content = string(data)
	return result, nil

}

type ActionProvider struct {
	Actions map[string]pvd.Action
}

func (p *ActionProvider) Init() {
	p.Actions = make(map[string]pvd.Action)
	p.Actions["read"] = pvd.Action{ActionFunc: ReadFile, Inputs: ReadFileInput}
}

func (p *ActionProvider) GetActions() map[string]pvd.Action {
	return p.Actions
}

func New() pvd.ActionProvider {
	return &ActionProvider{}
}
