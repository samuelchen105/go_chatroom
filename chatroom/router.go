package chatroom

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yuhsuan105/go_chatroom/common"
)

func HandlerRegister(rt *mux.Router) {
	rt.HandleFunc("/", showAll).Methods("GET")

}

func HandlerRegisterWithAuth(rt *mux.Router) {
	rt.Use(common.AuthHandler)
	rt.HandleFunc("/user", showUserChat).Methods("GET")
	//rt.HandleFunc("/create", showCreate).Methods("GET")
}

func showAll(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page != "" {

		val := common.GetSession(r, "showAll")
		chatrooms, ok := val.([]Chatroom)

		if !ok {
			log.Println("GetSession: something wrong")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		index, err := strconv.Atoi(page)
		if err != nil {
			log.Println("atoi: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pagenum := math.Ceil(float64(len(chatrooms)) / 10.0)
		selectArr := make([]int, int(pagenum))
		for i := range selectArr {
			selectArr[i] = i + 1
		}

		var showlen int
		if index*10 < len(chatrooms) {
			showlen = index * 10
		} else {
			showlen = len(chatrooms)
		}
		data := &templData{
			Chatrooms: chatrooms[(index-1)*10 : showlen],
			Select:    selectArr,
			Prev:      selectArr[0] - 1,
			Next:      selectArr[len(selectArr)-1] + 1,
		}

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

	err = common.SetSession(w, r, "showAll", chatrooms)
	if err != nil {
		log.Println("session: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chatrooms/?page=1", http.StatusFound)
}

func showUserChat(w http.ResponseWriter, r *http.Request) {

}

func Create(w http.ResponseWriter, r *http.Request) {

}
