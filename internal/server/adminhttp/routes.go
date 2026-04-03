package adminhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codelif/hostbin/internal/protocol/adminv1"
	"github.com/codelif/hostbin/internal/server/middleware"
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
	authenticated.GET(adminv1.AuthCheckRelativePath, handler.AuthCheck)
	authenticated.GET(adminv1.DocumentsRelativePath, handler.ListDocuments)
	authenticated.GET(adminv1.DocumentPathPattern, handler.GetDocument)
	authenticated.GET(adminv1.DocumentContentPattern, handler.GetDocumentContent)
	authenticated.POST(adminv1.DocumentPathPattern, handler.CreateDocument)
	authenticated.PUT(adminv1.DocumentPathPattern, handler.ReplaceDocument)
	authenticated.DELETE(adminv1.DocumentPathPattern, handler.DeleteDocument)

	engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, adminv1.ErrorResponse{Error: adminv1.ErrorNotFound})
	})
	engine.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, adminv1.ErrorResponse{Error: adminv1.ErrorMethodNotAllowed})
	})

	return engine
}
