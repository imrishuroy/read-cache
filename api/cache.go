package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	db "github.com/imrishuroy/read-cache-api/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func (server *Server) ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
	//ctx.JSON(http.StatusOK, "Ok")
}

type createCacheRequest struct {
	Title string `json:"title" binding:"required"`
	Link  string `json:"link" binding:"required"`
}

func (server *Server) createCache(ctx *gin.Context) {
	var req createCacheRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	// add this cache to DB
	arg := db.CreateCacheParams{
		Owner: authPayload.UID,
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
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if cache.Owner != authPayload.UID {
		err := errors.New("cache does't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	arg := db.ListCachesParams{
		Owner:  authPayload.UID,
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

type updateCacheRequest struct {
	ID       int64  `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Link     string `json:"link" binding:"required"`
	IsPublic bool   `json:"is_public"`
}

func (server *Server) updateCache(ctx *gin.Context) {
	var req updateCacheRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	dbCache, err := server.store.GetCache(context.Background(), req.ID)
	if err != nil {

		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no data found for this cache id"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if dbCache.Owner != authPayload.UID {
		err := errors.New("cache does't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateCacheParams{
		ID:       req.ID,
		Title:    req.Title,
		Link:     req.Link,
		IsPublic: pgtype.Bool{Bool: req.IsPublic, Valid: true},
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	dbCache, err := server.store.GetCache(context.Background(), req.ID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no data found for this cache id"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if dbCache.Owner != authPayload.UID {
		err := errors.New("cache does't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	error := server.store.DeleteCache(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
}

type listPublicCachesByTagIDsRequest struct {
	TagIDs []int32 `form:"tag_ids" binding:"required"`
}

func (server *Server) listPublicCaches(ctx *gin.Context) {
	var req listPublicCachesByTagIDsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("req.TagIDs", req.TagIDs)

	caches, err := server.store.ListPublicCaches(ctx, req.TagIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, caches)
}
