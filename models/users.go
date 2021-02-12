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

// userValidator type
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// UserGorm type
type UserGorm struct {
	db *gorm.DB
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
	hmac := hash.NewHMAC(hashSecretKey)
	uv := &userValidator{
		UserDB: ug,
		hmac:   hmac}
	return &userService{
		uv,
	}, nil
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

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create will create the given user
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.setRememberTokenIfUnset,
		uv.hmacRemember)

	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) setRememberTokenIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pBytes := []byte(user.Password + passwordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}

	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// Update will update given user info
func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the given user
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

// newUserGorm contructor to create new gorm
func newUserGorm(connString string) (*UserGorm, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &UserGorm{
		db: db,
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
func (ug *UserGorm) ByRemember(tokenHash string) (*User, error) {
	var u User
	db := ug.db.Where("remember_hash = ?", tokenHash)
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

// Create will create the given user
func (ug *UserGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will update given user info
func (ug *UserGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete the given user
func (ug *UserGorm) Delete(id uint) error {
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
