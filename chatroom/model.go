package chatroom

import (
	"errors"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/yuhsuan105/go_chatroom/common"
)

const (
	CHATROOM_NAME_MINLEN = 4
)

type Chatroom struct {
	Id         uint `gorm:"primaryKey"`
	Name       string
	Owner_name string
	Created_on string
	Chats      string
}

type templChatroom struct {
	Chatrooms []Chatroom
	Select    []int
	Prev      int
	Current   int
	Next      int
}

type templForm struct {
	CsrfField template.HTML
	ErrMsg    string
}

func chatroomNameValidate(name string, minLen int) error {
	r := []rune(name)
	if len(r) < minLen {
		return errors.New("chatroom name length too short")
	}

	var count int64

	db := common.GetDatabase()
	db.Table("chatrooms").Where("name = ?", name).Count(&count)

	if count > 0 {
		return errors.New("this name have been used")
	}
	return nil
}

func listPager(r *http.Request, key string, param string) (*templChatroom, error) {
	val := common.GetSession(r, key)
	list, ok := val.([]Chatroom)
	if !ok {
		log.Println("GetSession: something wrong")
		return nil, errors.New("GetSession: something wrong")
	}

	if len(list) == 0 {
		return &templChatroom{
			Chatrooms: nil,
			Select:    []int{1},
			Prev:      1,
			Current:   1,
			Next:      1,
		}, nil
	}

	index, err := strconv.Atoi(param)
	if err != nil {
		log.Println("atoi: ", err)
		return nil, err
	}

	totalPage := math.Ceil(float64(len(list)) / 10.0)
	pageNums := make([]int, int(totalPage))
	for i := range pageNums {
		pageNums[i] = i + 1
	}

	indexLimit := common.Min(index*10, len(list))
	prev := common.Max(pageNums[0], index-1)
	next := common.Min(pageNums[len(pageNums)-1], index+1)

	return &templChatroom{
		Chatrooms: list[(index-1)*10 : indexLimit],
		Select:    pageNums,
		Prev:      prev,
		Current:   index,
		Next:      next,
	}, nil
}
