package common

import "html/template"

type TemplForm struct {
	CsrfField template.HTML
	ErrMsg    string
}
