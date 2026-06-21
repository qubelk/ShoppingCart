package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"user/internal/service"
	"user/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	serv *service.UserService
}

func New(serv *service.UserService) *UserHandler {
	return &UserHandler{
		serv: serv,
	}
}

func respondError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrEmailAlredyExist):
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "email already register",
		})
	case errors.Is(err, user.ErrInvalidEmail):
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid email",
		})
	case errors.Is(err, user.ErrTooShortPassword):
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "password is to short",
		})
	case errors.Is(err, user.ErrWeakPassword):
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "password is to weak",
		})
	case errors.Is(err, user.ErrUserNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req user.RegisterRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	u, err := h.serv.Register(ctx, &req)
	if err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	ctx.JSON(http.StatusCreated, u)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req user.LoginRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	u, err := h.serv.Login(ctx, &req)
	if err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	token, err := h.serv.GenerateToken(u.User.ID.String())
	if err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	ctx.SetCookie(
		"jwt-token",
		token,
		3600,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	id, exists := ctx.Get("id")
	if !exists {
		err := errors.New("invalid user id")
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	u, err := h.serv.GetProfile(ctx, userID)
	if err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, u)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	var req user.DeleteRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	if err := h.serv.Delete(ctx, &req); err != nil {
		respondError(ctx, err)
		user.LogError(err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
