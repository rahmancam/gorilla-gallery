package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var tpls map[string]*template.Template

func getTemplate(templ ...string) *template.Template {
	templs := []string{"templates/index.gohtml", "templates/layouts/base.gohtml", "templates/layouts/footer.gohtml"}
	templs = append(templs, templ...)
	return template.Must(template.ParseFiles(templs...))
}

func init() {
	tpls = make(map[string]*template.Template)
	tpls["contact"] = getTemplate("templates/contact.gohtml")
	tpls["home"] = getTemplate("templates/home.gohtml")
	tpls["faq"] = getTemplate("templates/faq.gohtml")
	tpls["404"] = getTemplate("templates/404.gohtml")

}

func render(templ string, w http.ResponseWriter, data interface{}) {
	tpl, found := tpls[templ]

	if !found {
		panic("template not found: " + templ)
	}

	if err := tpl.Execute(w, data); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	render("contact", w, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	render("home", w, nil)
}

func faq(w http.ResponseWriter, r *http.Request) {
	render("faq", w, nil)
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render("404", w, nil)
}

func main() {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(pageNotFound)
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/faq", faq)
	http.ListenAndServe(":8080", router)
}
