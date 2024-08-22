package utils

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(ctx *gin.Context) (page int64, limit int64, err error) {
	limitParam := strings.TrimSpace(ctx.Query("limit"))
	if limitParam == "" {
		limitParam = "10"
	}

	limit, err = strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		log.Printf("failed to parse limit query param: %v", err)
		return page, limit, errors.New("failed to parse limit query param")
	}

	if limit == 0 {
		limit = 10
	}

	pageParam := strings.TrimSpace(ctx.Query("page"))
	if pageParam == "" {
		pageParam = "1"
	}

	page, err = strconv.ParseInt(pageParam, 10, 64)
	if err != nil {
		log.Printf("failed to parse page query param: %v", err)
		return page, limit, errors.New("failed to parse page query param")
	}

	if page == 0 {
		page = 1
	}

	return page, limit, nil
}
