package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/read-cache/db/sqlc"
)

type createCacheRequest struct {
	Title string `json:"title" binding:"required"`
	Link  string `json:"link" binding:"required"`
}

func (server *Server) createCache(ctx *gin.Context) {
	var req createCacheRequest
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
		errorCode := db.ErrorCode(err)
		if errorCode == db.ForeignKeyViolation || errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cache)
}

type getCacheRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCache(ctx *gin.Context) {
	var req getCacheRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cache, err := server.store.GetCache(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			//errorMessage := "Cache not found with the given ID " + strconv.FormatInt(req.ID, 10)
			//ctx.JSON(http.StatusNotFound, gin.H{"error": errorMessage})

			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cache)
}

type listCacheRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCaches(ctx *gin.Context) {
	var req listCacheRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCachesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	caches, err := server.store.ListCaches(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, caches)
}

type updateAccountRequest struct {
	ID    int64  `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Link  string `json:"link" binding:"required"`
}

func (server *Server) updateCache(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateCacheParams{
		ID:    req.ID,
		Title: req.Title,
		Link:  req.Link,
	}

	cache, err := server.store.UpdateCache(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cache)
}

type deleteCacheRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteCache(ctx *gin.Context) {
	var req deleteCacheRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteCache(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

}
