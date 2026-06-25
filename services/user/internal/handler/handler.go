package handler

import (
	"errors"
	"net/http"
	"pkg/logs"
	"user/internal/service"
	"user/internal/user"

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
	var req user.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid register request",
		})
		logs.LogError(err)
		return
	}

	u, err := h.serv.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": user.ErrEmailAlredyExist,
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusCreated, u)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req user.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid login request",
		})
		logs.LogError(err)
		return
	}

	u, err := h.serv.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "authorization failed",
		})
		logs.LogError(err)
		return
	}

	token, err := h.serv.GenerateToken(u.User.ID, u.User.Login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.SetCookie(
		"jwt-token",
		token,
		43200,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	login, ok := ctx.Get("login")
	if !ok {
		err := errors.New("invalid user login")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err,
		})
		logs.LogError(err)
		return
	}

	u, err := h.serv.GetProfile(ctx.Request.Context(), login.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user with this login not found",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, u)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	var req user.DeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid delete request",
		})
		logs.LogError(err)
		return
	}

	if err := h.serv.Delete(ctx.Request.Context(), &req); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user with this login not found",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
