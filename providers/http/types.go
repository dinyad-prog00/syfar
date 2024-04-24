package http

import (
	"net/url"
)

type Result struct {
	TimeSeconds float64                `json:"timeseconds,omitempty"`
	Status      string                 `json:"status,omitempty"`
	StatusCode  int                    `json:"statuscode,omitempty"`
	Request     HTTPRequest            `json:"request,omitempty"`
	Body        string                 `json:"body,omitempty"`
	Json        interface{}            `json:"json,omitempty"`
	Headers     map[string]interface{} `json:"headers,omitempty"`
	Err         string                 `json:"err,omitempty"`
	Systemout   string                 `json:"systemout,omitempty"`
}

type HTTPRequest struct {
	URL      string                 `json:"url,omitempty"`
	Method   string                 `json:"method,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
	Query    map[string]interface{} `json:"query,omitempty"`
	Body     interface{}            `json:"body,omitempty"`
	Headers  Headers                `json:"headers,omitempty"`
	Form     url.Values             `json:"form,omitempty"`
	PostForm url.Values             `json:"post_form,omitempty"`
}

type Headers map[string]string
