package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server *Server) Login(ctx *gin.Context) {

	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//ctx.JSON(http.StatusBadRequest, errorResponse(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Emai and password are required"})
		return
	}

	customToken, err := server.authService.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": customToken})
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server *Server) Register(ctx *gin.Context) {

	// DB OPS

	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//ctx.JSON(http.StatusBadRequest, errorResponse(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Emai and password are required"})
		return
	}

	customToken, err := server.authService.Register(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": customToken})

}
