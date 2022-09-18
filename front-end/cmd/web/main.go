package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Only / because we only have 1 page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {
	// Go template files to be parsed together with the given template
	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	// Add the page to the slice
	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	// Combine the templates with the page
	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	// Parse the page and templates
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
