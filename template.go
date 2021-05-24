package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	tmpl "text/template"

	log "github.com/sirupsen/logrus"
)

type Template struct {
	T          *tmpl.Template
	SearchPath []string
	names      map[string]string
}

// Get an instance of Template (text/template's Template with all funcs added)
func New(name string) *Template {
	return &Template{
		T:     tmpl.New(name).Funcs(Funcs).Option("missingkey=error"),
		names: make(map[string]string),
	}
}

// Load a template. It will search backward through the SearchPath.
func (t *Template) load(name string) error {
	if t.T.Lookup(name) != nil {
		return nil
	}
	if filepath.Base(name) != name {
		return fmt.Errorf("the template name should not include a path, use the search path: %s", name)
	}
	for i := len(t.SearchPath) - 1; i >= 0; i-- {

		fn := filepath.Join(t.SearchPath[i], name)
		_, err := t.T.ParseFiles(fn)
		if os.IsNotExist(err) { // try in next path
			log.Debugf("template not found: %s\n", fn)
			continue
		}
		if err != nil {
			return fmt.Errorf("could not load template %s: %s", fn, err)
		}
		log.Debugf("template loaded %d. %s %s\n", i, name, fn)
		t.names[name] = fn
		return nil
	}
	return fmt.Errorf("could not find template %s in search path", name)
}

func execute(tmpl *tmpl.Template, vars map[string]interface{}) (string, error) {
	varsP, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		varsP = []byte(fmt.Sprintf("%s", vars))
	}
	log.Debugf("execute template %v vars=%v\n", tmpl, varsP)
	var buf strings.Builder
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t *Template) Execute(vars map[string]interface{}) (string, error) {
	return execute(t.T, vars)
}

func (t *Template) ExecuteTemplate(name string, vars map[string]interface{}) (string, error) {
	err := t.load(name)
	if err != nil {
		return "", err
	}
	res, err := execute(t.T.Lookup(name), vars)
	if err != nil {
		return "", fmt.Errorf("could not render template %s: %s", t.names[name], err)
	}
	return res, nil
}
