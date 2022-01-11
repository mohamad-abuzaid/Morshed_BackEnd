package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is our User example model.
// Keep note that the tags for public-use (for our web app)
// should be kept in other file like "web/viewmodels/user.go"
// which could wrap by embedding the datamodels.User or
// define completely new fields instead but for the sake
// of the example, we will use this datamodel
// as the only one User model in our application.
type User struct {
	ID             int64     `db:"id" json:"id" form:"id"`
	Firstname      string    `db:"firstname" json:"firstname" form:"firstname"`
	Username       string    `db:"username" json:"username" form:"username"`
	HashedPassword []byte    `db:"_" json:"-" form:"-"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at" form:"updated_at"`
}

// TableName returns the database table name of a User.
func (u User) TableName() string {
	return "users"
}

// PrimaryKey returns the primary key of a User.
func (u *User) PrimaryKey() string {
	return "id"
}

// SortBy returns the column name that
// should be used as a fallback for sorting a set of User.
func (u *User) SortBy() string {
	return "updated_at"
}

// Scan binds mysql rows to this User.
func (u *User) Scan(rows *sql.Rows) error {
	u.CreatedAt = new(time.Time)
	u.UpdatedAt = new(time.Time)
	return rows.Scan(&u.ID, &u.Firstname, &u.Username, &u.CreatedAt, &u.UpdatedAt)
}

// Users is a list of products. Implements the `Scannable` interface.
type Users []*User

// Scan binds mysql rows to this Categories.
func (us *Users) Scan(rows *sql.Rows) (err error) {
	cp := *us
	for rows.Next() {
		u := new(User)
		if err = u.Scan(rows); err != nil {
			return
		}
		cp = append(cp, u)
	}

	if len(cp) == 0 {
		return sql.ErrNoRows
	}

	*us = cp

	return rows.Err()
}

// IsValid can do some very very simple "low-level" data validations.
func (u User) IsValid() bool {
	return u.ID > 0
}

// GeneratePassword will generate a hashed password for us based on the
// user's input.
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword will check if passwords are matched.
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
