package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var tpls map[string]*template.Template

func init() {
	tpls = make(map[string]*template.Template)
	tpls["contact"] = getTemplate("views/templates/contact.gohtml")
	tpls["home"] = getTemplate("views/templates/home.gohtml")
	tpls["faq"] = getTemplate("views/templates/faq.gohtml")
	tpls["signup"] = getTemplate("views/templates/signup.gohtml")
	tpls["404"] = getTemplate("views/templates/404.gohtml")
}

func getTemplate(templ ...string) *template.Template {
	templs := []string{"views/templates/index.gohtml"}
	templs = append(templs, getLayoutFiles()...)
	templs = append(templs, templ...)
	return template.Must(template.ParseFiles(templs...))
}

func getLayoutFiles() []string {
	files, err := filepath.Glob("views/templates/layouts/*.gohtml")
	if err != nil {
		panic(err)
	}
	return files
}

// Contact view
func Contact(w http.ResponseWriter, r *http.Request) {
	render("contact", w, nil)
}

// Home view
func Home(w http.ResponseWriter, r *http.Request) {
	render("home", w, nil)
}

// Faq view
func Faq(w http.ResponseWriter, r *http.Request) {
	render("faq", w, nil)
}

// Signup view
func Signup(w http.ResponseWriter, r *http.Request) {
	render("signup", w, nil)
}

// PageNotFound view
func PageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render("404", w, nil)
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
