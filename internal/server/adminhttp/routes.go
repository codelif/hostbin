package adminhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hostbin/internal/protocol/adminv1"
	"hostbin/internal/server/middleware"
)

func NewEngine(handler *Handler, maxDocSize int64, authMiddleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.HandleMethodNotAllowed = true
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	engine.RemoveExtraSlash = false
	engine.UseRawPath = true
	engine.UnescapePathValues = false

	engine.GET(adminv1.HealthPath, handler.Health)

	authenticated := engine.Group(adminv1.BasePath)
	authenticated.Use(middleware.LimitBodyBytes(maxDocSize), authMiddleware)
	authenticated.GET("/documents", handler.ListDocuments)
	authenticated.GET("/documents/:slug", handler.GetDocument)
	authenticated.GET("/documents/:slug/content", handler.GetDocumentContent)
	authenticated.PUT("/documents/:slug", handler.PutDocument)
	authenticated.DELETE("/documents/:slug", handler.DeleteDocument)

	engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, adminv1.ErrorResponse{Error: adminv1.ErrorNotFound})
	})
	engine.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, adminv1.ErrorResponse{Error: adminv1.ErrorMethodNotAllowed})
	})

	return engine
}
