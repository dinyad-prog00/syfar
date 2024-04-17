package runner

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syfar/parser"
	"syfar/providers"
	"syfar/reporters"
	"syfar/types"

	"github.com/alecthomas/participle/v2"
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

func (s Syfar) GetActionFunc(key string) (providers.ActionFunc, error) {
	keys := strings.Split(key, "_")
	if len(keys) != 2 {
		return nil, fmt.Errorf("Error getting action function, key should be provider_action")
	}
	p := s.actionsProviders[keys[0]]
	if p == nil {
		return nil, fmt.Errorf("Error getting action provider: %s", keys[0])
	}
	f := p.ActionsFuncs()[keys[1]]
	if f == nil {
		return nil, fmt.Errorf("Error getting action: %s", key)
	}
	return f, nil

}

func (s Syfar) Init() {
	for _, p := range s.actionsProviders {
		p.Init()
	}
}

func (s Syfar) Run(filedir string, filename string) {
	var ps = participle.MustBuild[parser.SyfarFile](participle.Unquote())
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Erreur lors de la lecture du fichier: %v", err))
	}
	ast, err := ps.ParseString(filename, string(content))

	if err != nil {
		panic(err)
	}

	fimport := GetFromImport(*ast, ps, filedir)
	ast.Entries = PrependManyToList(ast.Entries, fimport)

	ctx := context.Background()
	SetValueToContext(&ctx, "syfar.rootdir", parser.Value{String: &filedir})

	InitializeContext(&ctx, *ast)

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
				fmt.Println(err)
			} else {
				result = append(result, ri...)
			}

		case v.Stepper != nil:
			ri, _ := RunStepper(&ctx, s, *v.Stepper, i)
			result = append(result, ri...)

		case v.Print != nil:
			RunPrint(&ctx, *v.Print)
		}

	}

	if len(result) != 0 {
		fmt.Println("\n___________________________________________________________________")
		reporters.ConsoleReporter(result)
		fmt.Println("\n___________________________________________________________________")
	}
}
