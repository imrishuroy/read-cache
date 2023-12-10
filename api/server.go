package api

import "github.com/gin-gonic/gin"

//Server serves HTTP requests for our banking service.

type Server struct {
	router *gin.Engine
}

func NewServer() (*Server, error) {

	server := &Server{}

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
