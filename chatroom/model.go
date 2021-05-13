package chatroom

type Chatroom struct {
	Id         uint `gorm:"primaryKey"`
	Name       string
	Owner_name string
	Created_on string
	Chats      string
}

type templData struct {
	Chatrooms []Chatroom
	Select    []int
	Prev      int
	Next      int
}
