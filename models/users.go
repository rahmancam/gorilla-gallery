package models

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ErrNotFound defines error when the resource not found
	ErrNotFound  = errors.New("user: resource not found")
	ErrInvalidID = errors.New("user: ID provided is Invalid")
)

// UserService type
type UserService struct {
	db *gorm.DB
}

// User type
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
}

// NewUserService contructor to create new user service
func NewUserService(connString string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil
}

// ByID queries and retuns user by Id
func (us UserService) ByID(id uint) (*User, error) {
	var u User
	db := us.db.Where("id = ?", id)
	err := first(db, &u)
	return &u, err
}

// ByEmail queries and returns user by email
func (us UserService) ByEmail(email string) (*User, error) {
	var u User
	db := us.db.Where("email = ?", email)
	err := first(db, &u)
	return &u, err
}

func first(db *gorm.DB, u *User) error {
	err := db.First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

// Create will create the given user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Update will update given user info
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete will delete the given user
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	u := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&u).Error
}

// ResetTables create and migrate
func (us *UserService) ResetTables() error {
	if err := us.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate migrate the table automatically
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
