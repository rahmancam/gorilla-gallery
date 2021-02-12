package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ErrNotFound defines error when the resource not found
	ErrNotFound = errors.New("user: resource not found")
	// ErrInvalidID defines error of invalid user id
	ErrInvalidID = errors.New("user: ID provided is Invalid")
	// ErrInvalidPassword defines error of invalid password
	ErrInvalidPassword = errors.New("user: incorrect password")
)

const passwordPepper = "BCTX&^591"

// UserService type
type UserService struct {
	db *gorm.DB
}

// User type
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
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

// Authenticate used to authenticate user
func (us UserService) Authenticate(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password+passwordPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return user, nil
}

// Create will create the given user
func (us *UserService) Create(user *User) error {
	pBytes := []byte(user.Password + passwordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
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
