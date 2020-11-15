package main

import (
	"log"
	"os"

	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
	"gorm.io/driver/sqlite"
)

type Database struct {
	users     map[string]NewUserRequest
	comments  map[int64]NewCommentRequest
	reactions map[int64]ReactionRequest
	engine *gorm.DB
}

func (db *Database) AddUser(nu *NewUserRequest) (err error) {
	user := User{
	    Name: nu.Name,
	    Email: nu.Email,
	    Password: nu.Password,
	}
	result := db.engine.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	nu.ID = user.Id
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

type User struct {
	Id int64
	Name string `gorm:"notnull"`
	Email string `gorm:"unique; notnull"`
	Password string `gorm:"notnull"`
	Active bool `gorm:"notnull; default:false"`
	Created int64 `gorm:"autoCreateTime"`
}

func GetDB() *Database {
	engine, err := gorm.Open(sqlite.Open(DBFILE), &gorm.Config{
		Logger: gl.Default.LogMode(gl.Warn),
	})
	if err != nil {
		log.Fatalf("Could not open connection db exiting.")
		return nil
	}
	err = engine.AutoMigrate(&User{})

	if err != nil {
		log.Fatalf("Could not create table %s", err.Error())
		return nil
	}
	if db == nil {
		db = &Database{
			users:     map[string]NewUserRequest{},
			comments:  map[int64]NewCommentRequest{},
			reactions: map[int64]ReactionRequest{},
			engine: engine,
		}
	}
	return db
}

const DBFILE="./test-db"

func NewDB() *Database {
	os.Remove(DBFILE)
	db = nil
	return GetDB()
}
