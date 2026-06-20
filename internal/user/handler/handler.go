package handler

import (
	"cart/internal/user/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	serv *service.UserService
}

func New(serv *service.UserService) *UserHandler {
	return &UserHandler{
		serv: serv,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
}
