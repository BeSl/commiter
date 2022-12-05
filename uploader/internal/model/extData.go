package model

import (
	"github.com/jmoiron/sqlx"
)

type (
	Users struct {
		Id       string `db:"id"`
		Extid    string `json:"uid"`
		Name     string `json:"name"`
		IsAdmin  bool
		GitLogin string
	}

	ExtDataProcessors struct {
		UidBase          string `json:"uid"`
		Name             string `json:"name"`
		DateEvents       string `json:"dateevent"`
		BinaryData       string `json:"dataProc"`
		Filename         string `json:"filename"`
		Processed        bool
		ErrorDescription string
		Expansion        string `json:"exp"`
	}
)

func NewUsers() *Users {
	return &Users{}
}

func NewExtDataProcessors() *ExtDataProcessors {
	return &ExtDataProcessors{}
}

func (u *Users) AddUser(db *sqlx.DB) error {

	us := NewUsers()
	txtQuery := "Select id from users where extid=$1"
	err := db.Get(&us, txtQuery, u.Id)
	if err != nil {
		return nil
	}

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (extId, name) VALUES ($1, $2)", u.Id, u.Name)
	tx.Commit()

	return nil
}

func CreateVersion(exp *ExtDataProcessors, u *Users, db *sqlx.DB) error {

	txtQuery := "INSERT INTO extprocVersion (extId, name, BinaryData,Filename) VALUES ($1, $2, $3, $4)"
	tx := db.MustBegin()
	tx.MustExec(txtQuery, exp.UidBase, exp.Name, exp.BinaryData, exp.Filename)
	tx.Commit()
	return nil
}
