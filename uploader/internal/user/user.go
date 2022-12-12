package user

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Ref      string
	ExtId    string `json: uid`
	IsAdmin  bool
	FullName string
	
	DateCreate time.Time
}

func NewUser() *User {
	return &User{}
}

func (u *User) Create(db *sqlx.DB) error {
	return nil
}

func (u *User) Read(db *sqlx.DB) error {
	return nil
}

func (u *User) Update(db *sqlx.DB) error {
	return nil
}

func (u *User) Delete(db *sqlx.DB) error {
	return nil
}

func (u *User) UserTabName(db *sqlx.DB) string {
	return "users"
}
