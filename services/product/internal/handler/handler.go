package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"product/internal/product"
	"product/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	serv *service.ProductService
}

func respondError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, product.ErrInvalidTitle):
		fallthrough
	case errors.Is(err, product.ErrInvalidDescription):
		fallthrough
	case errors.Is(err, product.ErrInvalidPrice):
		ctx.JSON(http.StatusBadRequest, err)
	case errors.Is(err, product.ErrProductNotFound):
		fallthrough
	case errors.Is(err, product.ErrProductNotExists):
		ctx.JSON(http.StatusNotFound, err)
	default:
		ctx.JSON(http.StatusInternalServerError, "internal server error")
	}
}

func (ph *ProductHandler) Create(ctx *gin.Context) {
	var req product.CreateProductRequest

	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	res, err := ph.serv.Create(ctx, req)
	if err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (ph *ProductHandler) SearchProducts(ctx *gin.Context) {
	var req product.SearchProductRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
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
	var req product.GetProductRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	res, err := ph.serv.GetProduct(ctx, &req)
	if err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	ctx.JSON(http.StatusFound, res)
}

func (ph *ProductHandler) Delete(ctx *gin.Context) {
	var req product.DeleteProductRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	if err := ph.serv.Delete(ctx, &req); err != nil {
		respondError(ctx, err)
		product.LogError(err)
		return
	}

	ctx.Redirect(http.StatusNoContent, "/search")
}
