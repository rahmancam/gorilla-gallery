package controllers

import (
	"fmt"
	"net/http"

	"github.com/rahmancam/gorilla-gallery/helpers"
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

// Create new user on signup
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var formData SignupForm
	if err := helpers.ParseForm(r, &formData); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, formData.Email, formData.Password)
}
