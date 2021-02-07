package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func contact(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "contact.gohtml", nil); err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		panic(err)
	}
}

func faq(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "faq.gohtml", nil); err != nil {
		panic(err)
	}
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 - you are in wrong route</h1>")
}

func main() {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(pageNotFound)
	router.HandleFunc("/", index)
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/faq", faq)
	http.ListenAndServe(":8080", router)
}
