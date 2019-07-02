package model

import (
	"github.com/gorilla/websocket"
)

//Message 信息结构体
type Message struct {
	Type      string `json:"type,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Name      string `json:"name,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

//Client 客户端结构
type Client struct {
	ID     string
	Name   string
	Socket *websocket.Conn
	Send   chan []byte
}

//SetName 设置名称
func (c *Client) SetName(name string) string {
	c.Name = name
	return c.Name
}
