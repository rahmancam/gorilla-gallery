package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `Please contact us @ 
	<a href="mailto:helpdesk@gorilla-gallery.com">mailto:helpdesk@gorilla-gallery.com<a>`)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>Gorilla Gallery</h1>`)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/contact", contact)
	http.ListenAndServe(":8080", router)
}
