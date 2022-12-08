package model

import (
	"github.com/jmoiron/sqlx"
)

type (
	Users struct {
		Id       int64  `db:"id"`
		Extid    string `db:"extid" json:"extID"`
		Name     string `db:"name" json:"name"`
		IsAdmin  bool
		GitLogin string
	}

	ExtDataProcessors struct {
		UidBase          string `json:"uid"`
		Name             string `json:"name"`
		Base64data       string `json:"base64data"`
		Type             string `json:"type"`
		ExtID            string `json:"extID"`
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
	tx.MustExec(txtQuery, exp.UidBase, exp.Name, exp.Base64data, exp.Name)
	tx.Commit()
	return nil
}

var table_task = `CREATE TABLE public.commit_tasks (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	"name" varchar(200) NULL,
	extid uuid NULL,
	"type" varchar NULL,
	base64data text NULL,
	textcommit text NULL,
	dataevent date NULL,
	"user" bigint NULL,
	processed bool NULL,
	error text NULL,
	CONSTRAINT commit_tasks_pk PRIMARY KEY (id)
);
`
