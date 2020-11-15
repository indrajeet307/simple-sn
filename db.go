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
	user := Users{
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
	comment := Comments{
		FromUserId: nc.FromUser,
		ToUserId: nc.ToUser,
		Body: nc.Body,
	}
	result := db.engine.Create(&comment)
	if result.Error != nil {
	}
	nc.ID = comment.Id
}

func (db *Database) GetWallComments(uid int64) (ncr WallCommentsResponse) {
	var comments []Comments
	db.engine.Where("ToUserId = ?", uid).Find(&comments)
	resComments := []NewCommentRequest{}
	for _, comment := range comments {
		resComments = append(resComments,
			NewCommentRequest{
				ID: comment.Id,
				ToUser: comment.ToUserId,
				FromUser: comment.FromUserId,
				Body: comment.Body,
			})
	}
	ncr.Comments = resComments
	return
}

func (db *Database) AddCommentReaction(rr *ReactionRequest) {

	cr := CommentReactions{
		CommentId: rr.CommentID,
		ReactionId: rr.ReactionID,
		Count: 1,
	}
	db.engine.Create(&cr)
}

func (db *Database) GetCommentReactions(cid int64) (lr ListReactions) {
	reactions := []ReactionRequest{}

	var crs []CommentReactions

	db.engine.Where("cid = ?", cid).Find(&crs)

	for _, reaction := range crs {
		reactions = append(reactions, ReactionRequest{
			CommentID: reaction.CommentId,
			ReactionID: reaction.ReactionId,
		})
	}
	lr.Reactions = reactions
	return
}

var db *Database = nil

type Users struct {
	Id int64
	Name string `gorm:"notnull"`
	Email string `gorm:"unique; notnull"`
	Password string `gorm:"notnull"`
	Active bool `gorm:"notnull; default:false"`
	Created int64 `gorm:"autoCreateTime"`
}

type Reactions struct {
	Id int64
	Name string `gorm:"notnull; unique"`
}

type CommentReactions struct {

	CommentId int64`gorm:"notnull; column:cid"`
	Comment Comments `gorm:"notnull; foreginKey"`

	ReactionId int64`gorm:"notnull; column:rid"`
	Reaction Reactions `gorm:"notnull; foreginKey"`

	Count int64 `gorm:"default:0"`
}

type Comments struct {
	Id int64

	FromUserId int64 `gorm:"notnull; column:FromUserId"`
	FromUser Users `gorm:"foreginKey:FromUserId"`

	ToUserId int64 `gorm:"notnull; column:ToUserId"`
	ToUser Users `gorm:"foreginKey:ToUserId"`

	Body string `gorm:"notnull"`
	Deleted bool `gorm:"notnull; default:false"`
}

func GetDB() *Database {
	engine, err := gorm.Open(sqlite.Open(DBFILE), &gorm.Config{
		Logger: gl.Default.LogMode(gl.Warn),
	})
	if err != nil {
		log.Fatalf("Could not open connection db exiting.")
		return nil
	}
	err = engine.AutoMigrate(&Users{},&Comments{},&Reactions{},&CommentReactions{})

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
