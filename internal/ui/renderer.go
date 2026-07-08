package ui

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v5"
)

func seq(start, end int) []int {
	s := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		s = append(s, i)
	}
	return s
}

func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"seq": seq,
	}
}

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(pattern string) *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.New("").Funcs(TemplateFuncMap()).ParseGlob(pattern)),
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c *echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
