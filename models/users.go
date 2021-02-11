package models

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ErrNotFound defines error when the resource not found
	ErrNotFound = errors.New("user: resource not found")
)

// UserService type
type UserService struct {
	db *gorm.DB
}

// User type
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
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
	err := us.db.Where("id = ?", id).First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create will create the given user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// ResetTables create and migrate
func (us *UserService) ResetTables() {
	us.db.Migrator().DropTable(&User{})
	us.db.Create(&User{})
	us.db.AutoMigrate(&User{})
}
