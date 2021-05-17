package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/yuhsuan105/go_chatroom/common"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SetHandler(rt *mux.Router) {
	rt.HandleFunc("/register", showRegister).Methods("GET")
	rt.HandleFunc("/register", doRegister).Methods("POST")
	rt.HandleFunc("/login", showLogin).Methods("GET")
	rt.HandleFunc("/login", doLogin).Methods("POST")
}

func showRegister(w http.ResponseWriter, r *http.Request) {
	data := templForm{CsrfField: csrf.TemplateField(r), ErrMsg: ""}
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
		data := templForm{CsrfField: csrf.TemplateField(r), ErrMsg: err.Error()}
		common.GenerateHTML(w, data, "layout", "user_register")
		return
	}
	log.Println("validate fine")
	//hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//insert into database
	db := common.GetDatabase()
	user := &User{Name: form.Nickname, Email: form.Email, Password: hashedPwd}
	err = db.Create(user).Error
	if err != nil {
		log.Println("db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("insert into database fine")
	//redirect to login
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func showLogin(w http.ResponseWriter, r *http.Request) {
	data := templForm{CsrfField: csrf.TemplateField(r), ErrMsg: ""}
	common.GenerateHTML(w, data, "layout", "user_login")
}

func doLogin(w http.ResponseWriter, r *http.Request) {
	form := struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
		Token    string `schema:"auth_token"`
	}{}

	r.ParseForm()

	if err := schema.NewDecoder().Decode(&form, r.PostForm); err != nil {
		log.Println("schema decode:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//search database
	db := common.GetDatabase()
	var user User
	err := db.Table("users").Where("email=?", form.Email).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("login: ", err)
			data := templForm{CsrfField: csrf.TemplateField(r), ErrMsg: err.Error()}
			common.GenerateHTML(w, data, "layout", "user_login")
			return
		} else {
			log.Println("database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	//validate
	if user.Email != form.Email || bcrypt.CompareHashAndPassword(user.Password, []byte(form.Password)) != nil {
		log.Println("user enter wrong email or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//setcookie
	uc := &common.UserCookie{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
	err = common.SetCookie(w, uc)
	if err != nil {
		log.Println("set cookie: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//redirect
	http.Redirect(w, r, "/chatroom/user", http.StatusFound)
}
