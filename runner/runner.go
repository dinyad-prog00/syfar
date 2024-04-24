package runner

import (
	"context"
	"fmt"
	"os"
	as "syfar/assertions"
	t "syfar/parser"
	pvd "syfar/providers"
	rt "syfar/types"

	"github.com/alecthomas/participle/v2"
)

func ParseFile(filedir string, filename string) (*t.SyfarFile, error) {
	var ps = participle.MustBuild[t.SyfarFile](participle.Unquote())
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier: %v", err)
	}
	ast, err := ps.ParseString(filename, string(content))

	if err != nil {
		return nil, err
	}
	fimport := GetFromImport(*ast, ps, filedir)
	ast.Entries = PrependManyToList(ast.Entries, fimport)
	return ast, nil
}

func RunExpectationItem(ctx *context.Context, rctx *rt.ActionResultContext, item t.ExpectationItem, index int) (rt.ExpectationItemResult, error) {
	if item.Symbolic == nil {
		return rt.ExpectationItemResult{}, fmt.Errorf("symbolic should not be null")
	}
	val := GetValueFromContextOrResult(ctx, *rctx, item.Symbolic.Key)

	err := as.ValueCompare(val, GetValue(ctx, *item.Symbolic.Value), item.Symbolic.Opp)

	if err == nil {
		return rt.ExpectationItemResult{Id: index, Passed: true}, nil
	}

	return rt.ExpectationItemResult{Id: index, Passed: false, Message: err.Error()}, err

}

func RunExpectation(ctx *context.Context, rctx *rt.ActionResultContext, exp t.Expectation, index int) (rt.ExpectationResult, error) {
	passed := true
	result := []rt.ExpectationItemResult{}

	for i, item := range exp.Items {
		ri, _ := RunExpectationItem(ctx, rctx, *item, i+1)

		result = append(result, ri)

		if !ri.Passed {
			passed = false
		}
	}
	return rt.ExpectationResult{Id: index, Passed: passed, Items: result}, nil
}

func RunTest(ctx *context.Context, rctx *rt.ActionResultContext, test t.Test, index int) (rt.TestResult, error) {
	if test.Skipped {
		return rt.TestResult{Id: index, State: rt.StateSkipped, Expectations: []rt.ExpectationResult{}, Description: test.Description}, nil

	}
	passed := rt.StatePassed
	result := []rt.ExpectationResult{}

	for i, exp := range test.Expectations {
		ri, _ := RunExpectation(ctx, rctx, *exp, i+1)

		result = append(result, ri)

		if !ri.Passed {
			passed = rt.StateFailed
		}
	}
	return rt.TestResult{Id: index, State: passed, Expectations: result, Description: test.Description}, nil
}

func RunTestSet(ctx *context.Context, rctx *rt.ActionResultContext, set t.TestSet, index int) ([]rt.TestResult, error) {

	result := []rt.TestResult{}

	for _, test := range set.Tests {
		test.Description = fmt.Sprintf("%s [set] > %s", set.Description, test.Description)
		ri, _ := RunTest(ctx, rctx, *test, 0)
		result = append(result, ri)

	}
	return result, nil
}

func RunStepper(ctx *context.Context, s Syfar, steps t.Stepper, index int) ([]rt.TestResult, error) {

	result := []rt.TestResult{}

	for _, step := range steps.Steps {
		if step.Action != nil {
			v := fmt.Sprintf("%s [Step]", steps.Id)
			step.Action.Prefix = &v
			ri, _ := RunAction(ctx, s, *step.Action, 0)
			result = append(result, ri...)
		}
	}
	return result, nil
}

func RunAction(ctx *context.Context, s Syfar, action t.Action, index int) ([]rt.TestResult, error) {

	result := []rt.TestResult{}

	act, err := s.GetAction(action.Type)
	if err != nil {
		return nil, err
	}

	actfunc := act.ActionFunc

	params, testSets, tests, outs := FilterActionAttributes(action, true)
	valErr := ValidateAction(ctx, action, params, act.Inputs)
	if valErr != nil {
		return nil, valErr
	}
	jsonData, err := ActionParametersToStringJSON(ctx, params, act.Inputs)
	if err != nil {
		return nil, err
	}
	rst, err := actfunc(ctx, jsonData)
	if err != nil {
		return nil, err
	}

	rctx := rt.ActionResultContext{Result: rst}

	m, err := pvd.ProviderResultToMap(rst)

	if err == nil {
		rctx.MapResult = m
	}

	dm, err := DumpToMap(rst)
	if err == nil {
		rctx.DumpMapResult = dm
	}

	for _, out := range outs {
		RunOut(ctx, &rctx, action.Id, *out)
	}

	for _, ts := range testSets {
		rst, _ := RunTestSet(ctx, &rctx, *ts, 0)
		result = append(result, rst...)
	}

	for _, test := range tests {
		rst, _ := RunTest(ctx, &rctx, *test, 0)
		result = append(result, rst)
	}

	return result, nil
}

func RunOut(ctx *context.Context, rctx *rt.ActionResultContext, id string, out t.Out) error {

	for _, v := range out.Variables {
		val := GetValueFromContextOrResult(ctx, *rctx, v.Identifier)
		*ctx = context.WithValue(*ctx, contextKey(fmt.Sprintf("%s.%s", id, v.Name)), val)
	}

	return nil
}

func RunPrint(ctx *context.Context, print t.Print) {
	fmt.Printf("\x1b[30m%s\x1b[0m\n", print.Pos.String())
	for _, p := range print.Variables {
		fmt.Printf("  %v\n\n", JsonString(ctx, *p))
	}
}
