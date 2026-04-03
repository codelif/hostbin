// Copyright (c) 2026 Harsh Sharma <harsh@codelif.in>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// SPDX-License-Identifier: MIT

package adminhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codelif/hostbin/internal/protocol/adminv1"
	"github.com/codelif/hostbin/internal/server/middleware"
)

func NewEngine(handler *Handler, maxDocSize int64, middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.HandleMethodNotAllowed = true
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	engine.RemoveExtraSlash = false
	engine.UseRawPath = true
	engine.UnescapePathValues = false
	engine.Use(middleware.NoStore())
	engine.GET(adminv1.HealthPath, handler.Health)

	authenticated := engine.Group(adminv1.BasePath)
	authenticated.Use(middleware.LimitBodyBytes(maxDocSize))
	authenticated.Use(middlewares...)
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
