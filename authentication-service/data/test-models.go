package data

import (
	"database/sql"
	"time"
)

type TestPostgresRepository struct {
	DB *sql.DB
}

func NewTestPostgresRepository(db *sql.DB) *TestPostgresRepository {
	return &TestPostgresRepository{
		DB: db,
	}
}

// This one just leave it as dumb as possible because we are only testing the handlers, not the DB itself
func (u *TestPostgresRepository) GetAll() ([]*User, error) {
	users := []*User{}
	return users, nil
}

func (u *TestPostgresRepository) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		Email:     "me@here.com",
		FirstName: "First",
		LastName:  "Last",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

func (u *TestPostgresRepository) GetOne(id int) (*User, error) {
	user := User{
		ID:        1,
		Email:     "me@here.com",
		FirstName: "First",
		LastName:  "Last",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

func (u *TestPostgresRepository) Update(user User) error {
	return nil
}

func (u *TestPostgresRepository) DeleteByID(id int) error {
	return nil
}

func (u *TestPostgresRepository) Insert(user User) (int, error) {
	return 2, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *TestPostgresRepository) ResetPassword(password string, user User) error {
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *TestPostgresRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}
