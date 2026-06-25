package handler

import (
	"net/http"
	"pkg/logs"
	"product/internal/product"
	"product/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	serv *service.ProductService
}

func New(serv *service.ProductService) *ProductHandler {
	return &ProductHandler{
		serv: serv,
	}
}

func (ph *ProductHandler) Create(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user ID"})
		return
	}

	var req product.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid create product request",
		})
		logs.LogError(err)
		return
	}

	res, err := ph.serv.Create(ctx.Request.Context(), req, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (ph *ProductHandler) SearchProducts(ctx *gin.Context) {
	var req product.SearchProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid search product request",
		})
		logs.LogError(err)
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "title parameter is required",
		})
		return
	}

	res, err := ph.serv.SearchProducts(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		logs.LogError(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (ph *ProductHandler) GetProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product id",
		})
		logs.LogError(err)
		return
	}

	var req product.GetProductRequest
	req.ID = id

	res, err := ph.serv.GetProduct(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "product with this id not found",
		})
		return
	}

	ctx.JSON(http.StatusFound, res)
}

func (ph *ProductHandler) Delete(ctx *gin.Context) {
	userIDInterface, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product id",
		})
		logs.LogError(err)
		return
	}

	var req product.DeleteProductRequest
	req.ID = id
	req.OwnerID = userID

	if err := ph.serv.Delete(ctx.Request.Context(), &req); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "not found product to delete",
		})
		logs.LogError(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
