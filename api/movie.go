package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mide7/go/movie_scrapper/internal/services"
	"github.com/mide7/go/movie_scrapper/internal/utils"
)

type MovieApiHandler struct {
	ctx   context.Context
	mongo services.IMovie
}

func newMovieApiHandler(ctx context.Context, mongo services.IMovie) *MovieApiHandler {
	return &MovieApiHandler{
		ctx:   ctx,
		mongo: mongo,
	}
}

func (m *MovieApiHandler) GetAllHandler(ctx *gin.Context) {

	page, limit, err := utils.GetPaginationParams(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movies, result, err := m.mongo.GetAll(m.ctx, &services.GetAllParams{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": movies, "meta": gin.H{"total": result.Total, "page": result.Page, "limit": result.Limit}})
}
