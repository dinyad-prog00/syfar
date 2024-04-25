package runner

import (
	"context"
	"fmt"
	"strings"
	"syfar/parser"
	"syfar/providers"
	"syfar/reporters"
	"syfar/types"
)

type Syfar struct {
	actionsProviders map[string]providers.ActionProvider
	reporters        map[string]string
}

func NewSyfar() Syfar {
	sf := Syfar{}
	sf.actionsProviders = make(map[string]providers.ActionProvider)
	sf.reporters = map[string]string{}
	return sf
}

func (s Syfar) RegisterActionProvider(key string, p providers.ActionProvider) {
	s.actionsProviders[key] = p
}

func (s Syfar) GetAction(key string) (*providers.Action, error) {
	keys := strings.Split(key, "_")
	if len(keys) != 2 {
		return nil, fmt.Errorf("Error getting action function, key should be provider_action")
	}
	p := s.actionsProviders[keys[0]]
	if p == nil {
		return nil, fmt.Errorf("Error getting action provider: %s", keys[0])
	}
	act, ok := p.GetActions()[keys[1]]
	if !ok {
		return nil, fmt.Errorf("Error getting action: %s", key)
	}
	return &act, nil

}

func (s Syfar) Init() {
	for _, p := range s.actionsProviders {
		p.Init()
	}
}

func (s Syfar) Validate(filedir string, filename string) (*parser.SyfarFile, context.Context, error) {

	ast, err := ParseFile(filedir, filename)
	if err != nil {
		return nil, nil, err
	}

	err = ValidateSyfarFile(*ast)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	SetValueToContext(&ctx, "syfar.rootdir", parser.Value{String: &filedir})

	InitializeContext(&ctx, *ast)

	err = ValidateActions(&ctx, s, *ast)
	if err != nil {
		return nil, nil, err
	}

	return ast, ctx, nil
}

func (s Syfar) Run(ast *parser.SyfarFile, ctx context.Context) error {

	result := []types.TestResult{}

	for i, v := range ast.Entries {

		switch {
		case v.Test != nil:

			ri, _ := RunTest(&ctx, nil, *v.Test, i+1)
			result = append(result, ri)
		case v.TestSet != nil:
			ri, _ := RunTestSet(&ctx, nil, *v.TestSet, i)
			result = append(result, ri...)
		case v.Action != nil:
			ri, err := RunAction(&ctx, s, *v.Action, i)
			if err != nil {
				return err
			}
			result = append(result, ri...)

		case v.Stepper != nil:
			ri, _ := RunStepper(&ctx, s, *v.Stepper, i)
			result = append(result, ri...)

		case v.Print != nil:
			RunPrint(&ctx, *v.Print)
		}

	}

	syfarResult := BuildSyfarResult(result)

	if len(result) != 0 {

		reporters.ConsoleReporter(syfarResult)
		fmt.Println("\n___________________________________________________________________")
	}

	return nil
}
