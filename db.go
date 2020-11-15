package main

import (
	"errors"
	"fmt"
)

type Database struct {
	users     map[string]NewUserRequest
	comments  map[int64]NewCommentRequest
	reactions map[int64]ReactionRequest
}

func (db *Database) AddUser(nu *NewUserRequest) (err error) {
	numUsers := len(db.users)
	nu.ID = int64(numUsers)
	_, found := db.users[nu.Email]
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

func (db *Database) GetWallComments(uid int64) (ncr WallCommentsResponse) {
	comments := []NewCommentRequest{}

	for _, comment := range db.comments {
		if comment.ToUser == uid {
			comments = append(comments, comment)
		}
	}
	ncr.Comments = comments
	return
}

func (db *Database) AddCommentReaction(rr *ReactionRequest) {
	numReactions := len(db.reactions)
	db.reactions[int64(numReactions)] = *rr
}

func (db *Database) GetCommentReactions(cid int64) (lr ListReactions) {
	reactions := []ReactionRequest{}

	for _, reaction := range db.reactions {
		if reaction.CommentID == cid {
			reactions = append(reactions, reaction)
		}
	}
	lr.Reactions = reactions
	return
}

var db *Database = nil

func GetDB() *Database {
	if db == nil {
		db = &Database{
			users:     map[string]NewUserRequest{},
			comments:  map[int64]NewCommentRequest{},
			reactions: map[int64]ReactionRequest{},
		}
	}
	return db
}

func NewDB() *Database {
	db = nil
	return GetDB()
}
