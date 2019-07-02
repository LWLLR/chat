package main

import (
	"chat/server/websocket"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	fmt.Println("Starting application...")
	websocket.Init()
}
