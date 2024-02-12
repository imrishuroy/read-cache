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

	tag, err := server.store.CreateTag(ctx, tagName)
	if err != nil {
		errorCode := db.ErrorCode(err)

		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, gin.H{"error": req.TagName + " already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, tag)

}

func (server *Server) listTags(ctx *gin.Context) {
	tags, err := server.store.ListTags(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

type tagIDsRequest struct {
	TagIDs []int32 `json:"tag_ids"`
}

type cacheIDsRequest struct {
	CacheID int64 `uri:"cache_id"`
}

func (server *Server) addTagToCache(ctx *gin.Context) {

	var tagIDsRequest tagIDsRequest
	if err := ctx.ShouldBindJSON(&tagIDsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var cacheIDsRequest cacheIDsRequest
	if err := ctx.ShouldBindUri(&cacheIDsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, tagID := range tagIDsRequest.TagIDs {
		arg := db.AddTagToCacheParams{
			CacheID: cacheIDsRequest.CacheID,
			TagID:   tagID,
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
	}

	ctx.JSON(http.StatusOK, successResponse())
}

type listCacheTagsRequest struct {
	CacheID int64 `uri:"cache_id" binding:"required,min=1"`
}

func (server *Server) listCacheTags(ctx *gin.Context) {

	var req listCacheTagsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tags, err := server.store.ListCacheTags(ctx, req.CacheID)

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

type deleteTagRequest struct {
	TagID int32 `uri:"tag_id" binding:"required,min=1"`
}

func (server *Server) deleteTag(ctx *gin.Context) {
	var req deleteTagRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err1 := server.store.DeleteTagFromUserTagsTable(ctx, req.TagID)
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err1))
		return
	}

	err2 := server.store.DeleteTagFromCacheTagsTable(ctx, req.TagID)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err2))
		return
	}

	err3 := server.store.DeleteTagFromTagsTable(ctx, req.TagID)
	if err3 != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err3))
		return
	}

	ctx.JSON(http.StatusOK, successResponse())
}

type deleteCacheTagsRequest struct {
	CacheID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteCacheTags(ctx *gin.Context) {
	var req deleteCacheTagsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteCacheTag(ctx, req.CacheID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse())
}
