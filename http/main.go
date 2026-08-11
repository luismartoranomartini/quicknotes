package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func noteList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/templates/home.html")
	if err != nil {
		http.Error(w, "Aconteceu um erro ao executar a página", http.StatusInternalServerError)
	}
	tmpl.Execute(w, nil)
}

func noteView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Nota nao encotrada", http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("views/templates/noteView.html")
	if err != nil {
		http.Error(w, "Aconteceu um erro ao executar a página", http.StatusInternalServerError)
	}
	tmpl.Execute(w, id)
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
	PORT := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/", noteList)
	mux.HandleFunc("/note/view", noteView)
	mux.HandleFunc("/note/create", noteCreate)

	fmt.Printf("Servidor rodando na porta %s\n", PORT)
	http.ListenAndServe(PORT, mux)
}
