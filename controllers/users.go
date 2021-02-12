package controllers

import (
	"fmt"
	"net/http"

	"github.com/rahmancam/gorilla-gallery/helpers"
	"github.com/rahmancam/gorilla-gallery/models"
)

// Users type
type Users struct {
	us *models.UserService
}

// NewUserController constructor
func NewUserController(us *models.UserService) *Users {
	return &Users{
		us: us,
	}
}

// SignupForm type holds all user signup form submit data
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginForm type holds all user login form submit data
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create new user on signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var formData SignupForm
	if err := helpers.ParseForm(r, &formData); err != nil {
		panic(err)
	}
	usr := models.User{
		Name:     formData.Name,
		Email:    formData.Email,
		Password: formData.Password,
	}

	if err := u.us.Create(&usr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, formData)
}

// Login allows user to login into the system
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var formData LoginForm
	if err := helpers.ParseForm(r, &formData); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(formData.Email, formData.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalid email address")
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password provided")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
