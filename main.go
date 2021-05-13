package main

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/yuhsuan105/go_chatroom/common"
	"github.com/yuhsuan105/go_chatroom/user"
)

func main() {
	//init database
	common.InitDatabase()
	//init secure cookie
	common.InitCookie()
	//set up router
	rt := mux.NewRouter()
	rt.HandleFunc("/", hello).Methods("GET")
	/*
		srt := rt.PathPrefix("/chatroom").Subrouter()
		chatroom.HandlerRegister(srt)
	*/
	rtUser := rt.PathPrefix("/user").Subrouter()
	user.HandlerRegister(rtUser)
	//set up csrf
	CSRF := csrf.Protect(
		[]byte(`123456789zxcvbnm,./asdfghjkl;'qw`),
		csrf.FieldName("auth_token"),
		csrf.Secure(false),
	)
	//start server
	log.Println("server started")
	log.Fatal(http.ListenAndServe(":8080", CSRF(rt)))
}

func hello(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Login    string
		Register string
	}{
		Login:    "./user/login",
		Register: "./user/register",
	}

	common.GenerateHTML(w, data, "layout", "welcome")
}
