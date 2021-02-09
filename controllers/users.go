package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

// User type
type User struct{}

// NewUserController constructor
func NewUserController() *User {
	return &User{}
}

// SignupForm type holds all user signup form submit data
type SignupForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var decoder = schema.NewDecoder()

// Create new user on signup
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	var formData SignupForm
	if err := decoder.Decode(&formData, r.PostForm); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, formData.Email, formData.Password)
}
