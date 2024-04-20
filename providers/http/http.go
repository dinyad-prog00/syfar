package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"
	pvd "syfar/providers"
	t "syfar/types"
)

/**
Input config
*/

var RequestInput = []t.Input{
	{Name: "url", Type: reflect.String, Required: true},
	{Name: "method", Type: reflect.String, Required: false},
	{Name: "body", Type: reflect.Map, Required: false},
	{Name: "headers", Type: reflect.Map, Required: false},
}

// HTTP Action Provider
type ActionProvider struct {
	Actions map[string]pvd.Action
}

func (p *ActionProvider) Init() {
	p.Actions = make(map[string]pvd.Action)
	p.Actions["request"] = pvd.Action{ActionFunc: p.Request, Inputs: RequestInput}
}

func (p *ActionProvider) GetActions() map[string]pvd.Action {
	return p.Actions
}

func (p *ActionProvider) Request(ctx *context.Context, params interface{}) (interface{}, error) {
	paramString, ok := params.(string)
	if !ok {
		return nil, fmt.Errorf("params arg should be a string")
	}

	input, err := pvd.JsonParametersToProviderInputType[HTTPRequest](paramString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	httpP := HttpProvider{}
	return httpP.Do(*input)
}

type HttpProvider struct {
}

func (p *HttpProvider) Do(input HTTPRequest) (interface{}, error) {
	// Initialisation de la structure Result

	result := Result{
		Request: input,
	}
	body, err := json.Marshal(input.Body)
	if err != nil {
		result.Err = err.Error()
		return result, err
	}
	// Création de la requête HTTP
	req, err := http.NewRequest(input.Method, input.URL, bytes.NewBufferString(string(body)))
	if err != nil {
		result.Err = err.Error()
		return result, err
	}

	// Configuration des en-têtes de la requête
	for key, value := range input.Headers {
		req.Header.Set(key, value)
	}

	// Exécution de la requête HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Err = err.Error()
		return result, err
	}

	result.Status = resp.Status
	result.StatusCode = resp.StatusCode

	result.Headers = map[string]interface{}{}
	for key, value := range resp.Header {
		result.Headers[key] = value[0]
	}

	// Lecture du corps de la réponse
	if isBodySupported(resp) {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result.Err = err.Error()
			return result, err
		}

		result.Body = string(body)

		// Parse response body JSON and set it to Result.Json
		if isBodyJSONSupported(resp) {
			var m interface{}
			decoder := json.NewDecoder(strings.NewReader(string(body)))
			decoder.UseNumber()
			if err := decoder.Decode(&m); err == nil {
				result.Json = m

			} else {
				result.Json = nil
			}
		} else {
			result.Json = nil
		}
	}
	// Retour du résultat
	return result, nil
}

// From venom
func isBodySupported(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	return isContentTypeSupported(contentType)
}

func isContentTypeSupported(contentType string) bool {
	contentType = parseContentType(contentType)
	switch {
	case strings.HasSuffix(contentType, "+json"):
		return true
	case strings.HasPrefix(contentType, "image/"), strings.HasPrefix(contentType, "audio/"), strings.HasPrefix(contentType, "video/"),
		strings.HasPrefix(contentType, "font/"), strings.HasPrefix(contentType, "application/vnd."):
		return false
	case strings.HasPrefix(contentType, "application/"):
		x := strings.SplitN(contentType, "/", 2)[1]
		switch x {
		case "octet-stream", "x-abiword", "vnd.amazon.ebook", "x-bzip", "x-bzip2", "x-csh", "msword", "epub+zip", "java-archive", "ogg", "pdf",
			"x-rar-compressed", "rtf", "x-sh", "x-shockwave-flash", "x-tar", "zip", "x-7z-compressed":
			return false
		}
	case strings.Contains(contentType, "multipart/form-data"):
		return false
	}
	return true
}

func parseContentType(contentType string) string {
	parsed, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	return parsed
}

func isBodyJSONSupported(resp *http.Response) bool {
	contentType := parseContentType(resp.Header.Get("Content-Type"))
	return strings.Contains(contentType, "application/json") || strings.HasSuffix(contentType, "+json")
}
