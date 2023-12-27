package api

import (
	db "github.com/imrishuroy/read-cache/db/sqlc"
	"github.com/imrishuroy/read-cache/util"

	"github.com/gin-gonic/gin"
)

//Server serves HTTP requests for our banking service.

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{config: config, store: store}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/", server.ping)

	router.POST("/caches", server.createCache)
	// id is URI parameter
	router.GET("/caches/:id", server.getCache)
	// here page_id and page_size is query parameters
	router.GET("/caches", server.listCaches)
	router.PUT("/caches", server.updateCache)
	router.DELETE("/caches/:id", server.deleteCache)

	server.router = router

}

// Start HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
