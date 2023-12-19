package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/read-cache/db/sqlc"
)

type createCache struct {
	Title string `json:"title" binding:"required"`
	Link  string `json:"link" binding:"required"`
}

func (server *Server) createCache(ctx *gin.Context) {
	var req createCache
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// add this cache to DB
	arg := db.CreateCacheParams{
		Title: req.Title,
		Link:  req.Link,
	}

	cache, err := server.store.CreateCache(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cache)
}
