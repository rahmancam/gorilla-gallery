package main

import (
	"net/http"

	"github.com/rahmancam/gorilla-gallery/router"
)

func main() {
	router := router.AppRouter()
	http.ListenAndServe(":8080", router)
}
