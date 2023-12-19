package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/read-cache/db/sqlc"
)

//Server serves HTTP requests for our banking service.

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{store: store}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/caches", server.createCache)

	server.router = router

}

// Start HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
