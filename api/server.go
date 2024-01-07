package api

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	db "github.com/imrishuroy/read-cache-api/db/sqlc"
	"github.com/imrishuroy/read-cache-api/util"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Server serves HTTP requests
type Server struct {
	config util.Config
	store  db.Store
	auth   *auth.Client
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal().Msg("Failed to create Firebase app")
	}

	fmt.Println("fb connection done ", app)

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Firebase auth client")
	}

	server := &Server{config: config, store: store, auth: auth}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/", server.ping).Use(CORSMiddleware())

	authRoutes := router.Group("/").Use(authMiddleware(server.auth))

	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.POST("/users", server.createUser)

	// TODO: check why some methods are camel case and why ther are Capital
	authRoutes.POST("/caches", server.createCache)
	// id is URI parameter
	authRoutes.GET("/caches/:id", server.getCache)
	// here page_id and page_size is query parameters
	authRoutes.GET("/caches", server.listCaches)
	authRoutes.PUT("/caches", server.updateCache)
	authRoutes.DELETE("/caches/:id", server.deleteCache)

	server.router = router

}

// func corsMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		ctx.Header("Access-Control-Allow-Origin", "*")
// 		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 		ctx.Header("Access-Control-Allow-Headers", "Content-Type")

// 		ctx.Next()
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Start HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
