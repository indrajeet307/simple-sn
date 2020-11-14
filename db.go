package main

import (
	"fmt"
	"errors"
)


type Database struct {
	users map[string]NewUserRequest
	comments map[int64]NewCommentRequest
}

func (db *Database) AddUser(nu *NewUserRequest) (err error) {
	numUsers := len(db.users)
	nu.ID = int64(numUsers)
	_, found := db.users[nu.Email];
	if found {
		return errors.New(fmt.Sprintf("Email %s already registered", nu.Email))
	}
	db.users[nu.Email] = *nu
	return nil
}

var db *Database = nil

func GetDB () (*Database) {
	if db == nil {
		db = &Database{
			map[string]NewUserRequest{},
			comments: map[int64]NewCommentRequest{},
		}
	}
	return db
}

