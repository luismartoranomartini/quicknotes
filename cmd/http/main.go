package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// home
func noteList(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"views/templates/base.html",
		"views/templates/pages/home.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "aconteceu um erro", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", nil)
}

func noteView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "nota nao encotrada", http.StatusNotFound)
		return
	}
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/note-view.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "aconteceu um erro", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", id)
}

// formulário
func noteNew(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/note-new.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "aconteceu um erro", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", nil)
}

func noteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "Criando uma nova nota")
}

func main() {
	config := loadConfig()
	fmt.Printf("Servidor rodando na porta %s", config.ServerPort)

	mux := http.NewServeMux()

	staticHandler := http.FileServer(http.Dir("views/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	mux.HandleFunc("/", noteList)
	mux.HandleFunc("/note/view", noteView)
	mux.HandleFunc("/note/new", noteNew)
	mux.HandleFunc("/note/create", noteCreate)

	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux)
}
