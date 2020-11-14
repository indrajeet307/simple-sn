package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"

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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)

	r.HandleFunc("/users", AddNewUser).Methods("POST")

	log.Print("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
