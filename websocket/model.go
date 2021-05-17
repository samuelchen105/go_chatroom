package websocket

import (
	"encoding/json"
	"log"
)

type Message struct {
	IsText       bool   `json:"IsText"`
	IsRegister   bool   `json:"IsRegister"`
	IsUnregister bool   `json:"IsUnregister"`
	Data         string `json:"Data"`
}

func newTextMessage(data string) *Message {
	return &Message{
		IsText: true,
		Data:   data,
	}
}

func newRegisterMessage(data string) *Message {
	return &Message{
		IsRegister: true,
		Data:       data,
	}
}

func newUnregisterMessage(data string) *Message {
	return &Message{
		IsUnregister: true,
		Data:         data,
	}
}

func jsonMustMarshal(msg *Message) []byte {
	jm, err := json.Marshal(msg)
	if err != nil {
		log.Println("json marshal: ", err)
		panic(err)
	}
	return jm
}
