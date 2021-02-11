package models

import (
	"fmt"
	"testing"
	"time"
)

func testingUserService() (*UserService, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "demo"
		password = "demo"
		dbname   = "gallery_test"
	)

	pgconnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := NewUserService(pgconnString)
	if err != nil {
		return nil, err
	}
	// clear tables
	us.ResetTables()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}

	u := User{
		Name:  "Abdul Rahman",
		Email: "abdul@testing.com",
	}

	if err := us.Create(&u); err != nil {
		t.Fatal(err)
	}

	if u.ID == 0 {
		t.Errorf("Exxpected ID > 0, Received %d", u.ID)
	}

	if time.Since(u.CreatedAt) > time.Duration(time.Second*5) {
		t.Errorf("Expected CreatedAt to be recent. Received %s", u.CreatedAt)
	}

	if time.Since(u.UpdatedAt) > time.Duration(time.Second*5) {
		t.Errorf("Expected UpdatedAt to be recent. Received %s", u.CreatedAt)
	}

}
