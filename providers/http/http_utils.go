package http

import (
	"fmt"
	urlpkg "net/url"
	"regexp"
	"strings"
)

func extractParamsPlaceholders(input string) []string {
	re := regexp.MustCompile(`:(\w+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	placeholders := make([]string, len(matches))
	for i, match := range matches {
		placeholders[i] = match[1]
	}
	return placeholders
}

func buildUrl(url string, params map[string]interface{}, query map[string]interface{}) (string, error) {
	// Replace path parameters in URL
	paramNames := extractParamsPlaceholders(url)
	for _, n := range paramNames {
		param, ok := params[n]
		if !ok {
			return "", fmt.Errorf("params \"%s\" is not provided", n)
		}
		url = strings.ReplaceAll(url, ":"+n, fmt.Sprintf("%v", param))
	}

	// Create URL with query parameters
	u, err := urlpkg.Parse(url)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, value := range query {
		q.Set(key, fmt.Sprintf("%v", value))
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
