package handler

import (
	"cart/internal/cart"
	"cart/internal/service"
	"net/http"
	"pkg/logs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartHandler struct {
	serv *service.CartService
}

func New(serv *service.CartService) *CartHandler {
	return &CartHandler{
		serv: serv,
	}
}

func (ch *CartHandler) GetCart(ctx *gin.Context) {
	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		logs.LogError(err)
		return
	}

	req := &cart.GetCartRequest{
		UserID: userID,
	}

	res, err := ch.serv.GetCart(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *CartHandler) AddItem(ctx *gin.Context) {
	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var req cart.AddItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		logs.LogError(err)
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		logs.LogError(err)
		return
	}

	req.UserID = userID

	res, err := h.serv.AddItem(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *CartHandler) UpdateQuantity(ctx *gin.Context) {
	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var req cart.UpdateQuantityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		logs.LogError(err)
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		logs.LogError(err)
		return
	}

	req.UserID = userID

	res, err := h.serv.UpdateQuantity(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *CartHandler) RemoveItem(ctx *gin.Context) {
	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var req cart.RemoveItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		logs.LogError(err)
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		logs.LogError(err)
		return
	}

	req.UserID = userID

	res, err := h.serv.RemoveItem(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *CartHandler) CleanCart(ctx *gin.Context) {
	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		logs.LogError(err)
		return
	}

	req := &cart.ClearCartRequest{
		UserID: userID,
	}

	if err := h.serv.CleanCart(ctx.Request.Context(), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "cart successfully cleared",
	})
}

func (h *CartHandler) GetCartTTL(ctx *gin.Context) {
	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		logs.LogError(err)
		return
	}

	req := &cart.GetCartTTLRequest{
		UserID: userID,
	}

	res, err := h.serv.GetCartTTL(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
