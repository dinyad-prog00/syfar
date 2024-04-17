package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	pvd "syfar/providers"
	r "syfar/runner"
)

type FileActionProvider struct {
}

func ReadFile(ctx *context.Context, params interface{}) interface{} {
	paramString, ok := params.(string)
	if !ok {
		return nil
	}

	input, err := pvd.JsonParametersToProviderInputType[FileProviderInput](paramString)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rootDir := r.GetValueFromContext(*ctx, "syfar.rootdir").(string)
	filePath := filepath.Join(rootDir, input.Path)
	data, err := os.ReadFile(filePath)
	if err != nil {

		return nil

	}

	result := Result{}
	result.Content = string(data)
	return result

}

type ActionProvider struct {
	Actions map[string]pvd.ActionFunc
}

func (p *ActionProvider) Init() {
	p.Actions = make(map[string]pvd.ActionFunc)
	p.Actions["read"] = ReadFile
}

func (p *ActionProvider) ActionsFuncs() map[string]pvd.ActionFunc {
	return p.Actions
}

func New() pvd.ActionProvider {
	return &ActionProvider{}
}
