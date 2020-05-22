package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

var templatesFolder = "templates"

//templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(templatesFolder, t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	// root
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// we use an initialized type struct here instead of a func (note Handle vs HandleFunc)
	// and we can do this only because our type templateHandler implements
	// the serveHTTP method that all Handlers need

	// start the webserver
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}

}
