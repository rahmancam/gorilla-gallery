package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rahmancam/gorilla-gallery/controllers"
	"github.com/rahmancam/gorilla-gallery/views"
)

// AppRouter initializes app routes
func AppRouter() *mux.Router {
	userController := controllers.NewUserController()
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(views.PageNotFound)
	router.HandleFunc("/", views.Home).Methods("GET")
	router.HandleFunc("/contact", views.Contact).Methods("GET")
	router.HandleFunc("/faq", views.Faq).Methods("GET")
	router.HandleFunc("/signup", views.Signup).Methods("GET")

	router.HandleFunc("/signup", userController.Create).Methods("POST")

	return router
}
