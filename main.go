package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	rt := mux.NewRouter()
	rt.HandleFunc("/", hello).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", rt))
}

func hello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.Write([]byte("hello world"))
}
