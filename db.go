package main

import (
	"fmt"
	"errors"
)


type Database struct {
	users map[string]NewUserRequest
}

func (db *Database) AddUser (nu *NewUserRequest) (err error) {
	numUsers := len(db.users)
	nu.ID = numUsers
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
		}
	}
	return db
}

