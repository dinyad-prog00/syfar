package runner

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	t "syfar/parser"
	rt "syfar/types"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/fsamin/go-dump"
)

const (
	contextKeyAlreadySet = "Variable already set"

	varAlreadDeclared = "variable `%s`(%s) is already declared"
)

type ContextError struct {
	cause error
}

func (e *ContextError) Error() string {
	return e.cause.Error()
}

func NewContextError(format string, a ...interface{}) *ContextError {
	return &ContextError{cause: fmt.Errorf(format, a...)}
}

func contextVarError(ctx *context.Context, key string, pos lexer.Position) {
	if GetValueFromContext(*ctx, key) != nil {
		panic(fmt.Sprintf(varAlreadDeclared, key, pos.String()))
	}

}

func GetMapValue2(data interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	var value interface{} = data

	for _, k := range keys {

		v, ok := value.(map[string]interface{})
		if !ok {
			return nil
		}

		val, ok := v[k]
		if !ok {
			return nil
		}

		value = val
	}

	return value
}

func GetMapValue(data interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	var value interface{} = data
	for _, k := range keys {
		// Vérifier si la clé contient un indice de tableau
		if strings.Contains(k, "[") && strings.Contains(k, "]") {
			arrKey := strings.Split(k, "[")
			arrName := arrKey[0]
			idxStr := strings.TrimRight(arrKey[1], "]")
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				return nil
			}

			// Obtenir la valeur de la clé de tableau
			if arrName != "" {
				v, ok := value.(map[string]interface{})
				if !ok {
					return nil
				}
				arr, ok := v[arrName]
				if !ok {
					return nil
				}
				arrSlice, ok := arr.([]interface{})
				if !ok {
					return nil
				}
				if idx < 0 || idx >= len(arrSlice) {
					return nil
				}
				value = arrSlice[idx]
			} else {
				arrSlice, ok := value.([]interface{})
				if !ok {
					return nil
				}
				if idx < 0 || idx >= len(arrSlice) {
					return nil
				}
				value = arrSlice[idx]
			}
		} else {
			// Sinon, procéder comme avant

			v, ok := value.(map[string]interface{})
			if !ok {
				//println(value.(map[string]interface{})["0"])
				return nil
			}

			val, ok := v[k]
			if !ok {
				return nil
			}

			value = val
		}
	}

	return value
}

func ExtractInterpolationVariableNames(s string) []string {
	re := regexp.MustCompile(`\${([^}]+)}`)
	matches := re.FindAllStringSubmatch(s, -1)
	var names []string
	for _, match := range matches {
		names = append(names, match[1])
	}
	return names
}

func PrependToList[T any](list []T, element T) []T {
	newList := make([]T, len(list)+1)
	newList[0] = element
	copy(newList[1:], list)
	return newList
}

func PrependManyToList[T any](list []T, elements []T) []T {
	newList := make([]T, len(list)+len(elements))
	copy(newList[len(elements):], list)
	copy(newList, elements)
	return newList
}

func IsInStringList(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func IndexInStringList(target string, list []string) int {
	for i, item := range list {
		if item == target {
			return i
		}
	}
	return -1
}

func GetArgVarlue(name string, argNames []string, argValues []*t.Value) t.Value {
	for i, n := range argNames {
		if name == n {
			return *argValues[i]
		}
	}
	return t.Value{}
}

func CreateFolder(rootDir, folderName string) error {
	folderPath := filepath.Join(rootDir, folderName)
	err := os.MkdirAll(folderPath, 0700) // 0700 permission for folder (only owner has read, write, execute permission)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(folderPath, fileName string) (*os.File, error) {
	filePath := filepath.Join(folderPath, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600) // 0600 permission for file (only owner has read, write permission)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// func AppendToFile(file *os.File, data string) error {
// 	_, err := file.WriteString(data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func AppendToFile(filePath, data string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		return err
	}
	return nil
}

func WriteToFile(filePath, data string) error {
	err := os.WriteFile(filePath, []byte(data), 0600) // 0600 permission for file (only owner has read, write permission)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func AssignmentsToJSON(asgmts []t.Assignment) t.JSON {
	json := t.JSON{}
	for _, asgmt := range asgmts {
		json.Attributes = append(json.Attributes, &t.JSONAttribute{Name: asgmt.Name, Value: asgmt.Value})
	}
	return json
}

func AssignmentsFilter(asgmts []*t.Assignment, names []string) []t.Assignment {
	result := []t.Assignment{}
	for _, asgmt := range asgmts {
		if IsInStringList(asgmt.Name, names) {
			result = append(result, *asgmt)
		}
	}
	return result
}

func PrefixString(initial string, prefix *string, sep string) string {
	if prefix == nil {
		return initial
	}
	return fmt.Sprintf("%s%s%s", *prefix, sep, initial)
}

func FilterActionAttributes(action t.Action, setPrefix bool) ([]*t.Assignment, []*t.TestSet, []*t.Test, []*t.Out) {
	params := []*t.Assignment{}
	testSets := []*t.TestSet{}
	tests := []*t.Test{}
	outs := []*t.Out{}

	for _, attr := range action.Attributes {
		switch {
		case attr.Parameter != nil:
			params = append(params, attr.Parameter)
		case attr.TestSet != nil:
			if setPrefix {
				attr.TestSet.Description = PrefixString(fmt.Sprintf("%s: %s > %s", action.Type, action.Id, attr.TestSet.Description), action.Prefix, " > ")
			}
			testSets = append(testSets, attr.TestSet)
		case attr.Test != nil:
			if setPrefix {
				attr.Test.Description = PrefixString(fmt.Sprintf("%s: %s > %s", action.Type, action.Id, attr.Test.Description), action.Prefix, " > ")
			}
			tests = append(tests, attr.Test)
		case attr.Out != nil:
			outs = append(outs, attr.Out)
		}
	}
	return params, testSets, tests, outs
}

func ActionParametersToStringJSON(ctx *context.Context, params []*t.Assignment, inputs []rt.Input) (string, error) {
	result := make(map[string]interface{})
	for _, param := range params {
		result[param.Name] = GetValue(ctx, *param.Value)
	}
	json, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func ValueArrayToListString(array t.Value) []string {
	list := []string{}
	for _, v := range array.Array {
		list = append(list, *v.String)
	}
	return list
}

func GetTemplate(name string) (string, error) {
	return ReadFile(fmt.Sprintf("templates/%s.sf", name))
}

func EncryptString(key []byte, plaintext string) (string, error) {
	// Convert plaintext to byte slice
	plaintextBytes := []byte(plaintext)

	// Generate a new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new cipher block mode for AES-GCM encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext using AES-GCM
	ciphertext := gcm.Seal(nonce, nonce, plaintextBytes, nil)

	// Encode the ciphertext to base64 for easier storage and transmission
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func JsonString(ctx *context.Context, data t.Value) interface{} {
	value := GetValue(ctx, data)
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Array, reflect.Map:
		jsonData, err := json.MarshalIndent(value, "  ", "\t")
		if err != nil {
			return value
		}
		return string(jsonData)
	default:
		return value
	}
}

func GetPartAfterArray(s string) (string, bool) {
	re := regexp.MustCompile("^Array(.*)")
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 && matches[1] != "" {
		return matches[1], true
	} else if len(matches) > 1 {
		return "Id", true
	}
	return s, false
}

func IsJSONIntKey(key string) bool {
	// Expression régulière pour vérifier le format "json[int]"
	regex := regexp.MustCompile(`json\[\d+\]`)
	// Vérification si la clé correspond à l'expression régulière
	return regex.MatchString(key)
}

func DumpToMap(data interface{}) (map[string]interface{}, error) {
	////////////////
	e := dump.NewDefaultEncoder()
	// e.ExtraFields.Len = true
	//e.ExtraFields.Type = true
	// e.ExtraFields.DetailedStruct = true
	// e.ExtraFields.DetailedMap = true
	//e.ExtraFields.DetailedArray = true
	//e.Prefix = "req"
	e.ArrayJSONNotation = true
	e.ExtraFields.UseJSONTag = true
	//e.DisableTypePrefix = false
	e.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultLowerCaseFormatter()}

	return e.ToMap(data)
}

func ResultToR(identifier string) string {
	identifier = strings.ToLower(identifier)
	if strings.HasPrefix(identifier, "r.") {
		return "result" + identifier[1:]
	}
	return identifier
}

func Capitalize(value string) string {
	return strings.ToUpper(value[:1]) + strings.ToLower(value[1:])
}

func GetValueFromContextOrResult(ctx *context.Context, rctx *rt.ActionResultContext, indentifier string) interface{} {
	if rctx != nil {
		indentifier = ResultToR(indentifier)
		if rctx.DumpMapResult != nil && rctx.DumpMapResult[indentifier] != nil {
			return rctx.DumpMapResult[indentifier]
		}

		indentifier2 := Capitalize(strings.Replace(indentifier, "result.", "", 1))

		value := GetMapValue(rctx.MapResult, indentifier2)
		if value != nil {
			return value
		}
		value = GetMapValue(rctx.Result, indentifier2)
		if value != nil {
			return value
		}
	}
	return GetValueFromContext(*ctx, indentifier)
}

func BuildSyfarResult(testsResult []rt.TestResult) rt.SyfarResult {
	nbPassed := 0
	nbFailed := 0
	nbSkipped := 0
	for _, r := range testsResult {
		switch r.State {
		case rt.StatePassed:
			nbPassed++
		case rt.StateFailed:
			nbFailed++
		case rt.StateSkipped:
			nbSkipped++
		}
	}

	return rt.SyfarResult{TestsResult: testsResult, NbTestsPassed: nbPassed, NbTestsFailed: nbFailed, NbTestSkipped: nbSkipped}
}

var GET = "GET"
var POST = "POST"
