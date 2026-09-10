// Package render
package render

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"quicknotes/views"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"
)

type RenderTemplate struct {
	session *scs.SessionManager
}

func NewRender(session *scs.SessionManager) *RenderTemplate {
	return &RenderTemplate{session: session}
}

func getTemplatePageFiles(tmpl *template.Template, page string, useFS bool) (*template.Template, error) {
	if useFS {
		return tmpl.ParseFS(views.Files, "templates/base.html", "templates/pages/"+page)
	}
	files := []string{
		"views/templates/base.html",
	}
	files = append(files, "views/templates/pages/"+page)
	return tmpl.ParseFiles(files...)

}

func getTemplateMailFiles(mailTempl string, useFS bool) (*template.Template, error) {
	if useFS {
		return template.ParseFS(views.Files, "templates/mails/"+mailTempl)
	}
	return template.ParseFiles("views/templates/mails/" + mailTempl)
}

func (rt *RenderTemplate) RenderPage(w http.ResponseWriter, r *http.Request, status int, page string, data any) error {

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

	useFS := !strings.Contains(r.Host, "localhost")
	tmpl, err := getTemplatePageFiles(tmpl, page, useFS)
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

func (rt *RenderTemplate) RenderMailBody(r *http.Request, mailTempl string, data any) ([]byte, error) {
	useFS := !strings.Contains(r.Host, "localhost")
	tmpl, err := getTemplateMailFiles(mailTempl, useFS)
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
