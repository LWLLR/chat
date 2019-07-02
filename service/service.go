package service

import (
	"bytes"
	"chat/model"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

//Service chat service
type Service struct {
	Clients    map[*model.Client]bool
	Broadcast  chan []byte
	Register   chan *model.Client
	Unregister chan *model.Client
}

//Dispatch 分析协议号调度
func (service *Service) Dispatch(c *model.Client, message []byte) {
	bf := bytes.NewBuffer(message[:2])

	var no uint16 //协议号
	data := message[2:]
	binary.Read(bf, binary.BigEndian, &no)
	fmt.Println(no)
	switch no {
	case 1:
		service.Login(c, data)
	case 2:
		service.sendMsg(c, data)
	}

}

//Start 开启客户端管理器
func (service *Service) Start() {
	for {
		select {
		case conn := <-service.Register:
			service.Clients[conn] = true
			jsonMessage, _ := json.Marshal(&model.Message{Type: "notification", Content: "一名用户已进入聊天室"})
			service.send(jsonMessage, conn)
		case conn := <-service.Unregister:
			if _, ok := service.Clients[conn]; ok {
				close(conn.Send)
				delete(service.Clients, conn)
				jsonMessage, _ := json.Marshal(&model.Message{Type: "notification", Content: "一名用户已离开聊天室"})
				service.send(jsonMessage, conn)
			}
		case message := <-service.Broadcast:
			for conn := range service.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(service.Clients, conn)
				}
			}

		}

	}
}

//Login 设置名称
func (service *Service) Login(c *model.Client, data []byte) {
	msgBf := bytes.NewBuffer(data)
	message := make([]rune, len(data)/4)
	binary.Read(msgBf, binary.BigEndian, &message)
	name := ""
	for _, i := range message {
		name += string(i)
	}
	c.SetName(name)
	alertMessage := "设置名称：\"" + name + "\"成功!"
	jsonMessage, _ := json.Marshal(&model.Message{Type: "success", Content: alertMessage})
	c.Send <- jsonMessage
}

func (service *Service) sendMsg(c *model.Client, data []byte) {
	msgBf := bytes.NewBuffer(data)
	message := make([]rune, len(data)/4)
	binary.Read(msgBf, binary.BigEndian, &message)
	if c.Name == "" {
		alertMessage := "请先输入名称！"
		jsonMessage, _ := json.Marshal(&model.Message{Type: "info", Content: alertMessage})
		c.Send <- jsonMessage
	} else {
		msg := ""
		for _, i := range message {
			msg += string(i)
		}
		jsonMessage, _ := json.Marshal(&model.Message{Type: "broadcast", Sender: c.ID, Name: c.Name, Content: msg})
		service.Broadcast <- jsonMessage
	}

}

func (service *Service) send(message []byte, ignore *model.Client) {
	for conn := range service.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

func (service *Service) Read(c *model.Client) {
	defer func() {
		service.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			service.Unregister <- c
			c.Socket.Close()
			break
		}
		service.Dispatch(c, message)

	}
}

func (service *Service) Write(c *model.Client) {
	defer func() {
		c.Socket.Close()
	}()
	welcome, _ := json.Marshal(&model.Message{Type: "success", Content: "连接成功！"})
	c.Socket.WriteMessage(websocket.TextMessage, welcome)
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
