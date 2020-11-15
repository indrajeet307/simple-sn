package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const HuddleGreeting = "Hello Huddle!!"

type Response struct {
    Message string `json:"message"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, s string){
	er := ErrorResponse{s}
	erj, _ := json.Marshal(er)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(string(erj)))
}

func sendJsonResponse(w http.ResponseWriter, v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to Marshal output")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(j)))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sendJsonResponse(w, Response{HuddleGreeting})
}

type NewUserRequest struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Active string `json:"-"`
}

type NewUserResponse struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Active bool `json:"active"`
}

type NewCommentRequest struct {
	ID int64 `json:"id"`
	ToUser int64 `json:"to_user"`
	FromUser int64 `json:"from_user"`
	Body string `json:"body"`
	ParentID int64 `json:"parent_id"`
}

type NewCommentResponse struct {
	ID int64 `json:"id"`
}

func readRequest(r *http.Request, v interface{}) (err error) {
	body := []byte{}
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	return err
}

func AddNewUser(w http.ResponseWriter, r *http.Request) {
	newUser := NewUserRequest{}
	err := readRequest(r, &newUser)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error reading request object %s", err.Error()))
		return
	}
	db = GetDB()
	err = db.AddUser(&newUser)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJsonResponse(w, NewUserResponse{newUser.ID, newUser.Name, newUser.Email, true})
}

func AddToWall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Invalid user ID"))
	}
	newComment := NewCommentRequest{}
	err = readRequest(r, &newComment)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error reading request object %s", err.Error()))
		return
	}
	newComment.ToUser = uid
	db = GetDB()
	db.AddComment(&newComment)
	sendJsonResponse(w, NewCommentResponse{newComment.ID})
}

type WallCommentsResponse struct {
	Comments []NewCommentRequest
}

func GetUserWall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Invalid user ID"))
	}
	db = GetDB()
	comments := db.GetWallComments(uid)
	sendJsonResponse(w, comments)
}

func AddCommentReply(w http.ResponseWriter, r *http.Request) {
	newComment := NewCommentRequest{}
	err := readRequest(r, &newComment)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error reading request object %s", err.Error()))
		return
	}
	db = GetDB()
	db.AddComment(&newComment)
	sendJsonResponse(w, NewCommentResponse{newComment.ID})
}

type ReactionRequest struct {
	CommentID int64 `json:"comment_id"`
	ReactionID int64 `json:"reaction_id"`
}

type ReactionResponse struct {
	CommentID int64 `json:"comment_id"`
}

func AddCommentReaction(w http.ResponseWriter, r *http.Request) {
	reaction := ReactionRequest{}
	err := readRequest(r, &reaction)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error reading request object %s", err.Error()))
		return
	}
	db = GetDB()
	db.AddCommentReaction(&reaction)
	sendJsonResponse(w, ReactionResponse{reaction.CommentID})
}

type ListReactions struct {
	Reactions []ReactionRequest `json:"reactions"`
}

func GetCommentReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentID"]
	cid, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Invalid comment ID %s %s", commentID, err.Error()))
		return
	}
	db = GetDB()
	listReactions := db.GetCommentReactions(cid)
	sendJsonResponse(w, listReactions)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)

	r.HandleFunc("/users", AddNewUser).Methods("POST")
	//r.HandleFunc("/users/{userID}", AddNewUser).Methods("DELETE")

	r.HandleFunc("/wall/{userID}", GetUserWall).Methods("GET")
	r.HandleFunc("/wall/{userID}", AddToWall).Methods("POST")

	r.HandleFunc("/comments", AddCommentReply).Methods("POST")
	//r.HandleFunc("/comments/{commentID}", GetComment).Methods("GET")
	//r.HandleFunc("/comments/{commentID}", DeleteComment).Methods("DELETE")

	r.HandleFunc("/reactions", AddCommentReaction).Methods("POST")
	r.HandleFunc("/reactions/{commentID}", GetCommentReaction).Methods("GET")

	log.Print("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
