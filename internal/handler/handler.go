package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	HandleWebSocket(c *gin.Context)
}
