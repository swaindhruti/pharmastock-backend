package ui

import (
	"html/template"
	"io"
	"path/filepath"

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
	pages     map[string]*template.Template
	shared    *template.Template
}

func NewTemplateRenderer(templatesDir string) *TemplateRenderer {
	shared := template.New("").Funcs(TemplateFuncMap())
	template.Must(shared.ParseFiles(
		filepath.Join(templatesDir, "layout.gohtml"),
		filepath.Join(templatesDir, "partials.gohtml"),
	))

	pageNames := []string{"login", "dashboard", "stockists", "retailers", "medicines", "inventory", "upload"}
	pages := make(map[string]*template.Template, len(pageNames))
	for _, name := range pageNames {
		tmpl := template.Must(shared.Clone())
		template.Must(tmpl.ParseFiles(filepath.Join(templatesDir, name+".gohtml")))
		pages[name] = tmpl
	}

	return &TemplateRenderer{pages: pages, shared: shared}
}

func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	if tmpl, ok := t.pages[name]; ok {
		return tmpl.ExecuteTemplate(w, "layout", data)
	}
	return t.shared.ExecuteTemplate(w, name, data)
}
