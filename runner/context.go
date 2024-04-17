package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	t "syfar/parser"

	"github.com/alecthomas/participle/v2"
)

type contextKey string

// Fonction pour initialiser et retourner un nouveau contexte avec des valeurs declarées
func InitializeContext(ctx *context.Context, file t.SyfarFile) {

	for _, entry := range file.Entries {
		switch {
		case entry.Variable != nil:
			contextVarError(ctx, fmt.Sprintf("var.%s", entry.Variable.Name), entry.Variable.Pos)
			SetValueToContext(ctx, fmt.Sprintf("var.%s", entry.Variable.Name), *entry.Variable.Value)

		case entry.VarSet != nil:
			for _, v := range entry.VarSet.Variables {
				contextVarError(ctx, fmt.Sprintf("vars.%s.%s", entry.VarSet.Id, v.Name), v.Pos)

				SetValueToContext(ctx, fmt.Sprintf("vars.%s.%s", entry.VarSet.Id, v.Name), *v.Value)
			}

		case entry.MultiVariable != nil:
			for _, v := range entry.MultiVariable.Variables {
				contextVarError(ctx, fmt.Sprintf("var.%s", v.Name), v.Pos)

				SetValueToContext(ctx, fmt.Sprintf("var.%s", v.Name), *v.Value)
			}

		case entry.SecretSet != nil:
			for _, v := range entry.SecretSet.Variables {
				contextVarError(ctx, fmt.Sprintf("secrets.%s.%s", entry.SecretSet.Id, v.Name), v.Pos)

				SetValueToContext(ctx, fmt.Sprintf("secrets.%s.%s", entry.SecretSet.Id, v.Name), *v.Value)
			}
		}
	}

}

func GetFromImport(file t.SyfarFile, ps *participle.Parser[t.SyfarFile], filedir string) []*t.Entry {
	result := []*t.Entry{}
	for _, entry := range file.Entries {
		switch {
		case entry.Import != nil:
			for _, f := range entry.Import.Files {
				content, err := os.ReadFile(filepath.Join(filedir, f))
				if err != nil {
					panic(fmt.Sprintf("Erreur lors de la lecture du fichier: %v", err))
				}

				ast, err := ps.ParseString(f, string(content))

				if err != nil {
					panic(err)
				}
				if ast != nil {
					result = append(result, ast.Entries...)
				}
			}

		}
	}

	return result
}

func GetValueFromContext(ctx context.Context, key string) interface{} {
	value := ctx.Value(contextKey(key))
	if value != nil {
		return value
	}
	keys := strings.Split(key, ".")

	rkey := ""
	for i, k := range keys {
		if i == 0 {
			rkey = k
		} else {
			rkey = fmt.Sprintf("%s.%s", rkey, k)
		}

		value := ctx.Value(contextKey(rkey))
		if value != nil {
			val := GetMapValue(value, strings.Join(keys[i+1:], "."))
			if val != nil {
				return val
			}
		}
	}

	return nil
}

func SetValueToContext(ctx *context.Context, key string, value t.Value) {
	if GetValueFromContext(*ctx, key) != nil {
		panic("Error: " + contextKeyAlreadySet)
	}
	*ctx = context.WithValue(*ctx, contextKey(key), GetValue(ctx, value))

}

func GetValue(ctx *context.Context, v t.Value) interface{} {
	switch {
	case v.Boolean != nil:
		return *v.Boolean
	case v.Identifier != nil:
		return GetValueFromContext(*ctx, *v.Identifier)
	case v.String != nil:
		val := *v.String
		interpolationVars := ExtractInterpolationVariableNames(val)
		for _, v := range interpolationVars {
			val = strings.ReplaceAll(val, fmt.Sprintf("${%s}", v), GetValueFromContext(*ctx, v).(string))
		}
		return val
	case v.Number != nil:
		return *v.Number
	case v.Json != nil:
		return GetJSONValue(ctx, *v.Json)
	case v.Array != nil:
		return v.Array
	case v.Map != nil:
		return v.Map
	case v.Any != nil:
		return v.Any
	}
	return nil
}

func GetJSONValue(ctx *context.Context, jsonV t.JSON) interface{} {
	result := map[string]interface{}{}

	for _, v := range jsonV.Attributes {
		if v.Value.Json != nil {
			result[v.Name] = GetJSONValue(ctx, *v.Value.Json)
		} else {
			result[v.Name] = GetValue(ctx, *v.Value)
		}
	}
	return result
}
