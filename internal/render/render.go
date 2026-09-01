// Package render
package render

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"
)

type RenderTenplate struct {
	session *scs.SessionManager
}

func NewRender(session *scs.SessionManager) *RenderTenplate {
	return &RenderTenplate{session: session}
}

func (rt *RenderTenplate) RenderPage(w http.ResponseWriter, r *http.Request, status int, page string, data any) error {
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
		"isAuthenticated": func() bool {
			return rt.session.Exists(r.Context(), "userID")
		},
		"userEmail": func() string {
			return rt.session.GetString(r.Context(), "userEmail")
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

func (rt *RenderTenplate) RenderMailBody(mailTempl string, data any) ([]byte, error) {
	tmpl, err := template.ParseFiles("views/templates/mails/" + mailTempl)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	w := &bytes.Buffer{}
	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return w.Bytes(), nil
}
