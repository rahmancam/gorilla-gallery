package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rahmancam/gorilla-gallery/controllers"
	"github.com/rahmancam/gorilla-gallery/models"
	"github.com/rahmancam/gorilla-gallery/views"
)

// AppRouter initializes app routes
func AppRouter() *mux.Router {
	usrService := getUserService()
	userController := controllers.NewUserController(usrService)
	usrService.AutoMigrate()
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(views.PageNotFound)
	router.HandleFunc("/", views.Home).Methods("GET")
	router.HandleFunc("/contact", views.Contact).Methods("GET")
	router.HandleFunc("/faq", views.Faq).Methods("GET")
	router.HandleFunc("/signup", views.Signup).Methods("GET")
	router.HandleFunc("/login", views.Login).Methods("GET")
	router.HandleFunc("/cookietest", userController.CookieTest).Methods("GET")

	router.HandleFunc("/signup", userController.Create).Methods("POST")
	router.HandleFunc("/login", userController.Login).Methods("POST")

	return router
}

func getUserService() models.UserService {
	us, err := models.NewUserService(getConnectionString())
	if err != nil {
		panic(err)
	}
	return us
}

func getConnectionString() string {
	const (
		host     = "localhost"
		port     = 5432
		user     = "demo"
		password = "demo"
		dbname   = "gallery"
	)

	pgconnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return pgconnString
}
