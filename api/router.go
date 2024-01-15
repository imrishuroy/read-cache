package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/", server.ping).Use(CORSMiddleware())

	authRoutes := router.Group("/api").Use(authMiddleware(server.auth))

	// users
	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.POST("/users", server.createUser)

	// caches
	authRoutes.POST("/caches", server.createCache)
	// id is URI parameter
	authRoutes.GET("/caches/:id", server.getCache)
	// here page_id and page_size is query parameters
	authRoutes.GET("/caches", server.listCaches)
	authRoutes.PUT("/caches", server.updateCache)
	authRoutes.DELETE("/caches/:id", server.deleteCache)
	authRoutes.GET("caches/public", server.listPublicCaches)

	// tags
	authRoutes.POST("/tags", server.createTag)
	authRoutes.GET("/tags", server.listTags)
	authRoutes.POST("/caches/add-tag", server.addTagToCache)
	authRoutes.GET("/caches/:id/tags", server.listCacheTags)
	authRoutes.POST("/tags/:tag_id/subscribe", server.subscribeTag)
	authRoutes.DELETE("/tags/:tag_id/unsubscribe", server.unsubscribeTag)
	authRoutes.GET("/users/tags/subscriptions", server.listUserSubscriptions)

	server.router = router

}

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
