package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

func render(w http.ResponseWriter, r *http.Request, status int, page string, data any) error {
	files := []string{
		"views/templates/base.html",
	}
	files = append(files, "views/templates/pages/"+page)
	tmpl := template.New("").Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"csrfToken": func() string {
			return csrf.Token(r)
		},
	})
	tmpl, err := tmpl.ParseFiles(files...)
	if err != nil {
		return err
	}

	buff := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buff, "base", data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	buff.WriteTo(w)
	return nil
}
