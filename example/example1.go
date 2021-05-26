package main

import (
	"github.com/kellerza/template"

	log "github.com/sirupsen/logrus"
)

func main() {
	p := []string{"./", "../test_data"}
	t, _ := template.New("test", template.SearchPath(p...))
	vars := map[string]interface{}{
		"a": "a",
	}
	tname := "tst.tmpl"
	log.Infof("Rendering %s\n", tname)
	res, err := t.ExecuteTemplate(tname, vars)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("OK\n%s\n", res)
}
