package runner

import (
	"context"
	"fmt"
	"strings"
	t "syfar/parser"
	pvd "syfar/providers"
	rt "syfar/types"

	"github.com/alecthomas/participle/v2/lexer"
)

func JoinErrors(errs ...error) error {
	var errorMessages []string

	for _, err := range errs {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	if len(errorMessages) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errorMessages, "\n"))
}

func JoinString(list ...string) string {
	return strings.Join(list, ".")
}

func ValidateAction(ctx *context.Context, action t.Action, params []*t.Assignment, inputs []rt.Input) error {
	found := []string{}
	valErrors := []error{}
	for _, param := range params {
		found = append(found, param.Name)
		value := GetValue(ctx, *param.Value)
		err := pvd.ValidateInput(param.Name, value, inputs, param.Pos.String())
		if err != nil {
			valErrors = append(valErrors, err)
		}
	}

	for _, inp := range inputs {
		if !IsInStringList(inp.Name, found) && inp.Required {
			valErrors = append(valErrors, fmt.Errorf("error at %s; argument %s is required", action.Pos.String(), inp.Name))
		}
	}

	return JoinErrors(valErrors...)
}

func ValidateSyfarFile(file t.SyfarFile) error {
	declas := map[string]lexer.Position{}
	errorList := []error{}

	for _, entry := range file.Entries {
		switch {
		case entry.Action != nil:
			key := JoinString("actions.", entry.Action.Id)
			if val, ok := declas[key]; ok {
				errorList = append(errorList, fmt.Errorf("error at %s; action with id \"%s\" is already declared here: %s", entry.Action.Pos.String(), entry.Action.Id, val.String()))
			} else {
				declas[key] = entry.Action.Pos
			}
		case entry.Variable != nil:
			key := JoinString("var.", entry.Variable.Name)
			if val, ok := declas[key]; ok {
				errorList = append(errorList, fmt.Errorf("error at %s; variable \"%s\" is already declared here: %s", entry.Variable.Pos.String(), entry.Variable.Name, val.String()))
			} else {
				declas[key] = entry.Variable.Pos
			}
		case entry.MultiVariable != nil:
			for _, v := range entry.MultiVariable.Variables {
				key := JoinString("var.", v.Name)
				if val, ok := declas[key]; ok {
					errorList = append(errorList, fmt.Errorf("error at %s; variable \"%s\" is already declared here: %s", v.Pos.String(), v.Name, val.String()))
				} else {
					declas[key] = v.Pos
				}
			}

		case entry.VarSet != nil:
			key := JoinString("vars.", entry.VarSet.Id)
			if val, ok := declas[key]; ok {
				errorList = append(errorList, fmt.Errorf("error at %s; variable set  with id \"%s\" is already declared here: %s", entry.VarSet.Pos.String(), entry.VarSet.Id, val.String()))
			} else {
				declas[key] = entry.VarSet.Pos
				for _, v := range entry.VarSet.Variables {
					name := fmt.Sprintf("%s.%s", entry.VarSet.Id, v.Name)
					key := JoinString("vars.", name)

					if val, ok := declas[key]; ok {
						errorList = append(errorList, fmt.Errorf("error at %s; variable \"%s\" is already declared here: %s", v.Pos.String(), name, val.String()))
					} else {
						declas[key] = v.Pos
					}
				}
			}
		case entry.SecretSet != nil:
			key := JoinString("secrets.", entry.SecretSet.Id)
			if val, ok := declas[key]; ok {
				errorList = append(errorList, fmt.Errorf("error at %s; secrets set  with id \"%s\" is already declared here: %s", entry.SecretSet.Pos.String(), entry.SecretSet.Id, val.String()))
			} else {
				declas[key] = entry.SecretSet.Pos
				for _, v := range entry.SecretSet.Variables {
					name := fmt.Sprintf("%s.%s", entry.SecretSet.Id, v.Name)
					key := JoinString("secrets.", name)

					if val, ok := declas[key]; ok {
						errorList = append(errorList, fmt.Errorf("error at %s; secret \"%s\" is already declared here: %s", v.Pos.String(), name, val.String()))
					} else {
						declas[key] = v.Pos
					}
				}
			}

		}

	}

	return JoinErrors(errorList...)
}

func ValidateActions(ctx *context.Context, s Syfar, file t.SyfarFile) error {
	errorList := []error{}

	for _, entry := range file.Entries {
		switch {
		case entry.Action != nil:

			// Validate action arguments
			act, err := s.GetAction(entry.Action.Type)
			if err != nil {
				return err
			}

			params, _, _, _ := FilterActionAttributes(*entry.Action, false)
			valErr := ValidateAction(ctx, *entry.Action, params, act.Inputs)
			if valErr != nil {
				errorList = append(errorList, valErr)
			}
		}
	}

	return JoinErrors(errorList...)
}
