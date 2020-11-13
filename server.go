package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

const HuddleGreeting = "Hello Huddle!!"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(HuddleGreeting))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	log.Print("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
