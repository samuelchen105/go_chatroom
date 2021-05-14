package chatroom

import (
	"errors"

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
