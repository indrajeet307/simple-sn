package main

import (
	"encoding/json"
	"log"
	"net/http"

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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	log.Print("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
