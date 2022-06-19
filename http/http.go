package http

import "github.com/gomarkdown/markdown"

// New returns a new instace of Converter
func New() *Converter {
	return &Converter{}
}

// Converter is the Slack markdwn Converter implementation
type Converter struct {
}

// Format returns a unique name for the converter
func (converter *Converter) Format() string {
	return "http"
}

// Parse will parse the standard markdown and return the converted data
func (converter *Converter) Parse(markdwn []byte) ([]byte, error) {
	return markdown.ToHTML(markdwn, nil, nil), nil
}
