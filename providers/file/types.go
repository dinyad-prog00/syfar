package file

type FileProviderInput struct {
	Path string `json:"path,omitempty"`
}

type Result struct {
	Content string
}
