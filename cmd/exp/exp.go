package main

import (
	"html/template"
	"os"
	"time"
)

type TemplateData struct {
	Name string
	Age  int
}

func main() {
	t, err := template.ParseFiles("layout1.html", "home.html", "footer.html", "header.html")
	tempo := time.Now().Year()

	// fmt.Println(t.Name()) // hello.html
	// // tmpl = tmpl.Lookup("a")
	// fmt.Println(t.DefinedTemplates())
	if err != nil {
		panic(err)
	}

	// err = tmpl.Execute(os.Stdout, nil)
	err = t.ExecuteTemplate(os.Stdout, "layout1.html", tempo)
	if err != nil {
		panic(err)
	}
}
