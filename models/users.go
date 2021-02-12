package models

import (
	"errors"

	"github.com/rahmancam/gorilla-gallery/hash"
	"github.com/rahmancam/gorilla-gallery/rand"
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
const hashSecretKey = "XVBTING^&#36893BFD"

// UserService is a set of methods to query and alter user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

// userService type
type userService struct {
	UserDB
}

// UserValidator type
type UserValidator struct {
	UserDB
}

// UserGorm type
type UserGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// User type
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqueIndex"`
}

// UserDB interface contains common methods to query and alter
// single user
type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	AutoMigrate() error
	ResetTables() error
}

// NewUserService contructor to create new user service
func NewUserService(connString string) (UserService, error) {
	ug, err := newUserGorm(connString)
	if err != nil {
		return nil, err
	}
	uVal := &UserValidator{ug}
	return &userService{
		uVal,
	}, nil
}

// newUserGorm contructor to create new gorm
func newUserGorm(connString string) (*UserGorm, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hashSecretKey)
	return &UserGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

// ByID queries and retuns user by Id
func (ug *UserGorm) ByID(id uint) (*User, error) {
	var u User
	db := ug.db.Where("id = ?", id)
	err := first(db, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ByEmail queries and returns user by email
func (ug *UserGorm) ByEmail(email string) (*User, error) {
	var u User
	db := ug.db.Where("email = ?", email)
	err := first(db, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ByRemember queries and returns user by token
func (ug *UserGorm) ByRemember(token string) (*User, error) {
	var u User
	db := ug.db.Where("remember_hash = ?", ug.hmac.Hash(token))
	err := first(db, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func first(db *gorm.DB, u *User) error {
	err := db.First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

// Authenticate used to authenticate user
func (us *userService) Authenticate(email, password string) (*User, error) {
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
func (ug *UserGorm) Create(user *User) error {
	pBytes := []byte(user.Password + passwordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
}

// Update will update given user info
func (ug *UserGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

// Delete will delete the given user
func (ug *UserGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	u := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&u).Error
}

// ResetTables create and migrate
func (ug *UserGorm) ResetTables() error {
	if err := ug.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate migrate the table automatically
func (ug *UserGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
