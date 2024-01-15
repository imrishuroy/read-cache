package api

import (
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/read-cache-api/db/sqlc"
)

type createTagRequest struct {
	TagName string `json:"tag_name"`
}

func (server *Server) createTag(ctx *gin.Context) {
	var req createTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tagName := strings.ToLower(req.TagName)

	_, err := server.store.CreateTag(ctx, tagName)
	if err != nil {
		errorCode := db.ErrorCode(err)

		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, gin.H{"error": req.TagName + " already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, successResponse())

}

func (server *Server) listTags(ctx *gin.Context) {
	tags, err := server.store.ListTags(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

type addTagToCacheRequest struct {
	CacheID int64 `json:"cache_id"`
	TagID   int32 `json:"tag_id"`
}

func (server *Server) addTagToCache(ctx *gin.Context) {
	var req addTagToCacheRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddTagToCacheParams{
		CacheID: req.CacheID,
		TagID:   req.TagID,
	}

	_, err := server.store.AddTagToCache(ctx, arg)
	if err != nil {
		errorCode := db.ErrorCode(err)

		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "tag already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, successResponse())

}

type listCacheTagsRequest struct {
	cacheID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) listCacheTags(ctx *gin.Context) {

	var req listCacheTagsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tags, err := server.store.ListCacheTags(ctx, req.cacheID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

type subscribeTagRequest struct {
	TagID int32 `uri:"tag_id" binding:"required,min=1"`
}

func (server *Server) subscribeTag(ctx *gin.Context) {
	var req subscribeTagRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	arg := db.SubscribeTagParams{
		UserID: authPayload.UID,
		TagID:  req.TagID,
	}

	_, err := server.store.SubscribeTag(ctx, arg)
	if err != nil {
		errorCode := db.ErrorCode(err)

		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "already subscribed"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, successResponse())

}

type unsubscribeTagRequest struct {
	TagID int32 `uri:"tag_id" binding:"required,min=1"`
}

func (server *Server) unsubscribeTag(ctx *gin.Context) {
	var req unsubscribeTagRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	arg := db.UnsubscribeTagParams{
		UserID: authPayload.UID,
		TagID:  req.TagID,
	}

	err := server.store.UnsubscribeTag(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse())
}

func (server *Server) listUserSubscriptions(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	tags, err := server.store.ListUserSubscriptions(ctx, authPayload.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tags)
}
