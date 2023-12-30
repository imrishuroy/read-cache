package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/read-cache-api/db/sqlc"
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

type registerUserResponse struct {
	AccessToken string  `json:"access_token"`
	User        db.User `json:"user"`
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

	// authUser, err := server.authService.FireAuth.VerifyIDToken(context.Background(), "")
	// if err != nil{
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// }

	arg := db.CreateUserParams{
		// ID:    authUser.UID,
		ID:    "2",
		Email: req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	customToken, err := server.authService.Register(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := registerUserResponse{
		AccessToken: customToken,
		User:        user,
	}

	ctx.JSON(http.StatusOK, res)

}
