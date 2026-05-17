package handler

import (
	"net/http"

	"github.com/entertrans/backend-bogor.git/pkg/response"
	"github.com/gin-gonic/gin"
)

type PingHandler interface {
	Ping(c *gin.Context)
}

type pingHandler struct{}

func NewPingHandler() PingHandler {
	return &pingHandler{}
}

func (h *pingHandler) Ping(c *gin.Context) {
	response.SendResponse(c, http.StatusOK, "pong", "Server is running")
}
