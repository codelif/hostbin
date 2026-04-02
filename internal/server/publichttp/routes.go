package publichttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewEngine(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.HandleMethodNotAllowed = true
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	engine.RemoveExtraSlash = false

	engine.GET("/", handler.GetRoot)
	engine.HEAD("/", handler.GetRoot)
	engine.NoRoute(func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Header("X-Content-Type-Options", "nosniff")
		c.String(http.StatusNotFound, "not found\n")
	})
	engine.NoMethod(func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Header("X-Content-Type-Options", "nosniff")
		c.String(http.StatusMethodNotAllowed, "method not allowed\n")
	})

	return engine
}
