package main

import (
	"errors"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

type Database struct {
	engine *gorm.DB
}

func (db *Database) AddUser(nu *NewUserRequest) (err error) {
	user := Users{
		Name:     nu.Name,
		Email:    nu.Email,
		Password: nu.Password,
	}
	result := db.engine.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	nu.ID = user.Id
	return nil
}

func (db *Database) AddWallComment(nc *NewCommentRequest) (err error) {
	comment := Comments{
		FromUserId: nc.FromUser,
		ToUserId:   nc.ToUser,
		Body:       nc.Body,
	}
	result := db.engine.Create(&comment)
	if result.Error != nil {
		return result.Error
	}
	nc.ID = comment.Id
	return
}

func (db *Database) GetWallComments(uid int64) (ncr WallCommentsResponse) {
	var comments []Comments
	db.engine.Where("ToUserId = ?", uid).Find(&comments)
	resComments := []NewCommentRequest{}
	for _, comment := range comments {
		resComments = append(resComments,
			NewCommentRequest{
				ID:       comment.Id,
				ToUser:   comment.ToUserId,
				FromUser: comment.FromUserId,
				Body:     comment.Body,
			})
	}
	ncr.Comments = resComments
	return
}

func (db *Database) AddCommentReaction(cid int64, rr *CommentReactionRequest) (res *CommentReactionResponse, err error) {

	res = &CommentReactionResponse{}

	err = db.engine.Transaction(func(tx *gorm.DB) error {
		cr := CommentReactions{}
		result := tx.Where("cid = ? AND rid = ?", cid, rr.ReactionID).First(&cr)
		if result.RowsAffected == 0 {

			cr.CommentId = cid
			cr.ReactionId = rr.ReactionID
			cr.Count = 1

			result := tx.Create(&cr)
			if result.Error != nil {
				return result.Error
			}
		} else {
			cr.Count += 1
			tx.Model(&CommentReactions{}).Where("cid = ? AND rid = ?", cid, rr.ReactionID).Update("count", cr.Count)
		}

		res.CommentID = cr.CommentId
		res.ReactionID = cr.ReactionId
		res.Count = cr.Count
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *Database) GetCommentReactions(cid int64) (lr CommentListReactions, err error) {

	reactions := []CommentReactionResponse{}

	var crs []CommentReactions

	results := db.engine.Where("cid = ?", cid).Find(&crs)

	if results.Error != nil {
		return lr, results.Error
	}

	for _, reaction := range crs {
		reactions = append(reactions, CommentReactionResponse{
			CommentID:  reaction.CommentId,
			ReactionID: reaction.ReactionId,
			Count:      reaction.Count,
		})
	}
	lr.Reactions = reactions
	return
}

func (db *Database) CheckPassword(req *SignInRequest) error {
	var user Users

	results := db.engine.Where("email = ?", req.Email).First(&user)
	if results.Error != nil {
		return results.Error
	}
	if user.Password != req.Password {
		return errors.New("Password dont match")
	}
	return nil
}

func (db *Database) AddReaction(req *ReactionRequest) error {
	reaction := Reactions{
		Name: req.Name,
	}
	result := db.engine.Create(&reaction)
	if result.Error != nil {
		return result.Error
	}
	req.ID = reaction.Id
	return nil
}

func (db *Database) ListReaction() (lr ReactionResponse, err error) {
	lr = ReactionResponse{}
	reactions := []Reactions{}
	result := db.engine.Find(&reactions)
	if result.Error != nil {
		return lr, result.Error
	}
	respReactions := []ReactionRequest{}

	for _, r := range reactions {
		respReactions = append(respReactions, ReactionRequest{
			ID:   r.Id,
			Name: r.Name,
		})
	}
	lr.Reactions = respReactions
	return
}

var db *Database = nil

type Users struct {
	Id       int64
	Name     string `gorm:"notnull"`
	Email    string `gorm:"unique; notnull"`
	Password string `gorm:"notnull"`
	Active   bool   `gorm:"notnull; default:false"`
	Created  int64  `gorm:"autoCreateTime"`
}

type Reactions struct {
	Id   int64
	Name string `gorm:"notnull; unique"`
}

type CommentReactions struct {
	CommentId int64    `gorm:"notnull; column:cid; index:react_index"`
	Comment   Comments `gorm:"notnull; foreginKey"`

	ReactionId int64     `gorm:"notnull; column:rid; index:react_index"`
	Reaction   Reactions `gorm:"notnull; foreginKey"`

	Count int64 `gorm:"default:0"`
}

type Comments struct {
	Id int64

	FromUserId int64 `gorm:"notnull; column:FromUserId"`
	FromUser   Users `gorm:"foreginKey:FromUserId"`

	ToUserId int64 `gorm:"notnull; column:ToUserId"`
	ToUser   Users `gorm:"foreginKey:ToUserId"`

	Body    string `gorm:"notnull"`
	Deleted bool   `gorm:"notnull; default:false"`
}

func GetDB() *Database {
	engine, err := gorm.Open(sqlite.Open(DBFILE), &gorm.Config{
		Logger:                 gl.Default.LogMode(gl.Warn),
		SkipDefaultTransaction: false,
	})
	if err != nil {
		log.Fatalf("Could not open connection db exiting.")
		return nil
	}
	err = engine.AutoMigrate(&Users{}, &Comments{}, &Reactions{}, &CommentReactions{})

	if err != nil {
		log.Fatalf("Could not create table %s", err.Error())
		return nil
	}
	engine.Exec("PRAGMA foreign_keys = ON")
	if db == nil {
		db = &Database{
			engine: engine,
		}
	}
	return db
}

const DBFILE = "./test-db"

func NewDB() *Database {
	os.Remove(DBFILE)
	db = nil
	return GetDB()
}
