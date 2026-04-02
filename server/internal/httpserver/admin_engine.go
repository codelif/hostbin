package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	adminpkg "hostbin/internal/admin"
)

func NewAdminEngine(handler *adminpkg.Handler, maxDocSize int64, authMiddleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.HandleMethodNotAllowed = true
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	engine.RemoveExtraSlash = false
	engine.UseRawPath = true
	engine.UnescapePathValues = false

	engine.GET("/api/v1/health", handler.Health)

	authenticated := engine.Group("/api/v1")
	authenticated.Use(LimitBodyBytes(maxDocSize), authMiddleware)
	authenticated.GET("/documents", handler.ListDocuments)
	authenticated.GET("/documents/:slug", handler.GetDocument)
	authenticated.GET("/documents/:slug/content", handler.GetDocumentContent)
	authenticated.PUT("/documents/:slug", handler.PutDocument)
	authenticated.DELETE("/documents/:slug", handler.DeleteDocument)

	engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, adminpkg.ErrorResponse{Error: "not_found"})
	})
	engine.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, adminpkg.ErrorResponse{Error: "method_not_allowed"})
	})

	return engine
}
