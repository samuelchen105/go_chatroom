package chatroom

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/yuhsuan105/go_chatroom/common"
)

func SetHandler(rt *mux.Router) {
	rt.HandleFunc("/", showAll).Methods("GET")
}

func SetHandlerWithAuth(rt *mux.Router) {
	rt.Use(common.AuthHandler)
	rt.HandleFunc("/user", showUserChatroom).Methods("GET")
	rt.HandleFunc("/create", showCreate).Methods("GET")
	rt.HandleFunc("/create", doCreate).Methods("POST")
	rt.HandleFunc("/chat", showChat).Methods("GET")
}

func showAll(w http.ResponseWriter, r *http.Request) {
	sessionKey := "showAll"
	page := r.URL.Query().Get("page")
	if page != "" {
		//pager
		data, err := listPager(r, sessionKey, page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//show
		common.GenerateHTML(w, data, "layout", "chatroom_all")
		return
	}
	db := common.GetDatabase()
	var chatrooms []Chatroom
	err := db.Table("chatrooms").Order("created_on").Find(&chatrooms).Error
	if err != nil {
		log.Printf("database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = common.SetSession(w, r, sessionKey, chatrooms)
	if err != nil {
		log.Println("session: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chatrooms/?page=1", http.StatusFound)
}

func showUserChatroom(w http.ResponseWriter, r *http.Request) {
	sessionKey := "showUserChatroom"
	page := r.URL.Query().Get("page")
	if page != "" {
		//pager
		data, err := listPager(r, sessionKey, page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//show
		common.GenerateHTML(w, data, "layout", "chatroom_user")
		return
	}
	//get cookie
	userCookie, err := common.ReadCookie(w, r)
	if err != nil {
		log.Println("read cookie:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//get from db
	db := common.GetDatabase()
	var chatrooms []Chatroom
	err = db.Table("chatrooms").Where("owner_name=?", userCookie.Name).Order("created_on").Find(&chatrooms).Error
	if err != nil {
		log.Printf("database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//set session
	err = common.SetSession(w, r, sessionKey, chatrooms)
	if err != nil {
		log.Println("session: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//redirect
	http.Redirect(w, r, "/chatroom/user?page=1", http.StatusFound)
}

func showCreate(w http.ResponseWriter, r *http.Request) {
	data := templForm{CsrfField: csrf.TemplateField(r), ErrMsg: ""}
	common.GenerateHTML(w, data, "layout", "chatroom_create")
}

func doCreate(w http.ResponseWriter, r *http.Request) {
	//read post form
	form := struct {
		ChatroomName string `schema:"chatroom_name"`
		Token        string `schema:"auth_token"`
	}{}

	r.ParseForm()
	if err := schema.NewDecoder().Decode(&form, r.PostForm); err != nil {
		log.Println("schema decode:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//validate
	err := chatroomNameValidate(form.ChatroomName, CHATROOM_NAME_MINLEN)
	if err != nil {
		data := templForm{CsrfField: csrf.TemplateField(r), ErrMsg: err.Error()}
		common.GenerateHTML(w, data, "layout", "chatroom_create")
		return
	}
	//get user
	cookie, err := common.ReadCookie(w, r)
	if err != nil {
		log.Println("read cookie:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//db
	db := common.GetDatabase()
	dbData := Chatroom{Name: form.ChatroomName, Owner_name: cookie.Name}
	err = db.Table("chatrooms").Select("Name", "Owner_name").Create(&dbData).Error
	if err != nil {
		log.Println("db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/chatroom/user", http.StatusFound)
}

func showChat(w http.ResponseWriter, r *http.Request) {
	chatId := r.URL.Query().Get("id")

	db := common.GetDatabase()
	chatroom := Chatroom{}
	db.Table("chatrooms").Where("id=?", chatId).Take(&chatroom)

	userCookie, err := common.ReadCookie(w, r)
	if err != nil {
		log.Println("read cookie:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := struct {
		ChatName  string
		ChatOwner string
		ChatId    string
		UserName  string
	}{
		ChatName:  chatroom.Name,
		ChatOwner: chatroom.Owner_name,
		ChatId:    chatId,
		UserName:  userCookie.Name,
	}
	common.GenerateHTML(w, data, "layout", "chatroom_chat")
}
