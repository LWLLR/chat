package websocket

import (
	"chat/model"
	"chat/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var manager = &service.Service{
	Broadcast:  make(chan []byte),
	Register:   make(chan *model.Client),
	Unregister: make(chan *model.Client),
	Clients:    make(map[*model.Client]bool),
}

//Init 初始化
func Init() {
	go manager.Start()
	r := gin.Default()
	router(r)
	// r.Run("0.0.0.0:12345")
	r.RunTLS("0.0.0.0:12345", "/root/ssl/lwlgo.com.crt", "/root/ssl/lwlgo.com.key")
}

func router(r *gin.Engine) {
	r.GET("/ws", wsPage)
}

func wsPage(c *gin.Context) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		return
	}
	uuid, _ := uuid.NewV4()

	client := &model.Client{ID: uuid.String(), Socket: conn, Send: make(chan []byte)}

	manager.Register <- client

	go manager.Read(client)
	go manager.Write(client)
}
