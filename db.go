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
func (db *Database) AddComment(nc *NewCommentRequest) {
	numComments := len(db.comments)
	nc.ID = int64(numComments)
	db.comments[nc.ID] = *nc
}

func (db *Database) GetWallComments(uid int64) (ncr WallCommentsResponse){
	comments := []NewCommentRequest{}

	for _, comment := range db.comments {
		if comment.ToUser == uid {
			comments = append(comments, comment)
		}
	}
	ncr.Comments = comments
	return
}

var db *Database = nil

func GetDB () (*Database) {
	if db == nil {
		db = &Database{
			users: map[string]NewUserRequest{},
			comments: map[int64]NewCommentRequest{},
		}
	}
	return db
}


func NewDB () (*Database) {
	db = nil
	return GetDB()
}
