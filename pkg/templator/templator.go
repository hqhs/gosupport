package templator

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
)

// Templator parses web/templates and builds them according to directory structure
type Templator struct {
	Root      string
	Templates map[string]*template.Template
}

// NewTemplator initializes Templator with data
func NewTemplator(root string) (*Templator, error) {
	t := make(map[string]*template.Template)
	tDir := path.Join(root, "web", "templates")
	// NOTE: current template building algorithms is subject of change
	// and I skipped any existence checks because if something is missing,
	// throw error and dont finish initialization process.
	base := path.Join(tDir, "base.tmpl")
	baseDir := path.Join(tDir, "base")
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			t[info.Name()] = overlay([]string{ base, path })
		}
		return nil
	})
	// TODO add emails
	return &Templator{tDir, t}, err
}

// GetTemplates return slice of available templates
func (t *Templator) GetTemplates() []string {
	a := make([]string, 0)
	for k := range t.Templates {
		a = append(a, k)
	}
	return a
}

func overlay(components []string) *template.Template {
	return template.Must(template.ParseFiles(components...))
}
