package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createCache struct {
	Title string `json:"title" binding:"required"`
	Link  string `json:"link" binding:"required"`
}

type cacheResponse struct {
	Title     string    `json:"title" binding:"required"`
	Link      string    `json:"link" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createCache(ctx *gin.Context) {
	var req createCache
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// add this cache to DB

	res := cacheResponse{
		Title:     req.Title,
		Link:      req.Link,
		CreatedAt: time.Now(),
	}

	ctx.JSON(http.StatusOK, res)

}
