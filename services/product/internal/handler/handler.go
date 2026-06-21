package handler

import (
	"errors"
	"net/http"
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

func respondError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, product.ErrInvalidTitle):
		fallthrough
	case errors.Is(err, product.ErrInvalidDescription):
		fallthrough
	case errors.Is(err, product.ErrInvalidPrice):
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	case errors.Is(err, product.ErrProductNotFound):
		fallthrough
	case errors.Is(err, product.ErrProductNotExists):
		ctx.JSON(http.StatusNotFound, gin.H{"message": err})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	}
}

func (ph *ProductHandler) Create(ctx *gin.Context) {
	login, exists := ctx.Get("login")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	var req product.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	res, err := ph.serv.Create(ctx, req, login.(string))
	if err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (ph *ProductHandler) SearchProducts(ctx *gin.Context) {
	var req product.SearchProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "title parameter is required",
		})
		return
	}

	res, err := ph.serv.SearchProducts(ctx, &req)
	if err != nil {
		respondError(ctx, err)
		product.LogError(err)
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
		product.LogError(err)
		return
	}

	var req product.GetProductRequest
	req.ID = id

	res, err := ph.serv.GetProduct(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "product with this id not found",
		})
		return
	}

	ctx.JSON(http.StatusFound, res)
}

func (ph *ProductHandler) Delete(ctx *gin.Context) {
	login, exists := ctx.Get("login")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product id",
		})
		product.LogError(err)
		return
	}

	var req product.DeleteProductRequest
	req.ID = id
	req.Owner = login.(string)

	if err := ph.serv.Delete(ctx, &req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
