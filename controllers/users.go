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
	signIn(w, &usr)
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// Login allows user to login into the system
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var formData LoginForm
	if err := helpers.ParseForm(r, &formData); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(formData.Email, formData.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	signIn(w, user)
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func signIn(w http.ResponseWriter, user *models.User) {
	cookie := http.Cookie{
		Name:   "email",
		Value:  user.Email,
		Secure: false,
	}
	http.SetCookie(w, &cookie)
}

// CookieTest test
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Email is: ", cookie.Value)
}
