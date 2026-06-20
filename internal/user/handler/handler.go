package handler

import (
	"cart/internal/user"
	"cart/internal/user/service"
	"encoding/json"
	"errors"
	"net/http"

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

func respondError(ctx *gin.Context, code int, err error) {
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
			"message": "user not founded",
		})
	default:
		ctx.JSON(code, gin.H{
			"message": err.Error(),
		})
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req user.RegisterRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	u, err := h.serv.Register(ctx, &req)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, u)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req user.LoginRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	u, err := h.serv.Login(ctx, &req)
	if err != nil {
		respondError(ctx, http.StatusUnauthorized, err)
		return
	}

	token, err := h.serv.GenerateToken(u.User.ID.String())
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
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
		respondError(ctx, http.StatusBadRequest, errors.New("invalid user id"))
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	u, err := h.serv.GetProfile(ctx, userID)
	if err != nil {
		respondError(ctx, http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusOK, u)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	var req user.DeleteRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.serv.Delete(ctx, &req); err != nil {
		respondError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
