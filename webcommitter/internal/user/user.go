package user

import (
	"time"
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
