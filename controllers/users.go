package controllers

import "net/http"

// User type
type User struct{}

// NewUserController constructor
func NewUserController() *User {
	return &User{}
}

// Create new user on signup
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
}
