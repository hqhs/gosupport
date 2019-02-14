package templator

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

// Templator parses web/templates and builds them according to directory structure
type Templator struct {
	Root      string
	Templates map[string]*template.Template
	// unexported fields
	pool      *bpool.BufferPool
}

// NewTemplator initializes Templator with data
func NewTemplator(root string) (*Templator, error) {
	pool := bpool.NewBufferPool(64)
	t := make(map[string]*template.Template)
	tDir := path.Join(root, "web", "templates")
	// NOTE: current template building algorithms is subject of change
	// and I skipped any existence checks because if something is missing,
	// throw error and dont finish initialization process.
	base := path.Join(tDir, "base.tmpl")
	baseDir := path.Join(tDir, "base")
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			t[info.Name()] = overlay([]string{base, path})
		}
		return nil
	})
	// TODO add emails
	return &Templator{tDir, t, pool}, err
}

// Render renders data to template with name and writes it to w
func (t *Templator) Render(w http.ResponseWriter, name string, data map[string]interface{}) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return fmt.Errorf("Requested template is not built")
	}
	b := t.pool.Get()
	defer t.pool.Put(b)
	// NOTE here "base" is template from parsed files which we would render, better way
	// to do it is store parsed templates as set of graphs and use its root for render
	if err := tmpl.ExecuteTemplate(b, "base", data); err != nil {
		return err
	}
	w.Header().Set("content-type", "text/html; charset=utf-8")
	b.WriteTo(w)
	return nil
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
