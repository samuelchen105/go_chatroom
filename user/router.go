package user

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/yuhsuan105/go_chatroom/common"
)

func HandlerRegister(rt *mux.Router) {
	rt.HandleFunc("/register", showRegister).Methods("GET")
	rt.HandleFunc("/register", doRegister).Methods("POST")
	rt.HandleFunc("/login", login).Methods("POST")
}

func showRegister(w http.ResponseWriter, r *http.Request) {
	data := registerTempl{CsrfField: csrf.TemplateField(r), ErrMsg: ""}
	common.GenerateHTML(w, data, "layout", "user_register")
}

func doRegister(w http.ResponseWriter, r *http.Request) {
	form := struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
		Retype   string `schema:"retype"`
		Nickname string `schema:"nickname"`
		Token    string `schema:"auth_token"`
	}{}

	r.ParseForm()

	if err := schema.NewDecoder().Decode(&form, r.PostForm); err != nil {
		log.Println("schema decode:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("parse form fine")
	//validate
	err := registerValidate(form.Email, form.Password, form.Retype, form.Nickname)

	if err != nil {
		data := registerTempl{CsrfField: csrf.TemplateField(r), ErrMsg: err.Error()}
		common.GenerateHTML(w, data, "layout", "user_register")
		return
	}
	log.Println("validate fine")
	//insert into database
	db := common.GetDatabase()
	user := &User{Name: form.Nickname, Email: form.Email, Password: form.Password}
	err = db.Create(user).Error
	if err != nil {
		log.Println("db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("insert into database fine")
	//redirect to login
	common.Redirect(w, "/login")
}

func login(w http.ResponseWriter, r *http.Request) {

}
