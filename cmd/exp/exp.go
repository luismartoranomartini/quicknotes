package main

import (
	"html/template"
	"os"
)

type TemplateData struct {
	Name string
	Age  int
}

func main() {
	tmpl, err := template.ParseFiles("hello.html")
	// fmt.Println(tmp.Name()) // hello.html
	if err != nil {
		panic(err)
	}

	data := TemplateData{Name: "Luis", Age: 49}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
