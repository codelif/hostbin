package adminv1

import "net/url"

const (
	BasePath               = "/api/v1"
	HealthPath             = BasePath + "/health"
	AuthCheckPath          = BasePath + "/auth/check"
	DocumentsCollection    = BasePath + "/documents"
	AuthCheckRelativePath  = "/auth/check"
	DocumentsRelativePath  = "/documents"
	DocumentPathPattern    = "/documents/:slug"
	DocumentContentPattern = "/documents/:slug/content"
)

func DocumentPath(slug string) string {
	return DocumentsCollection + "/" + url.PathEscape(slug)
}

func DocumentContentPath(slug string) string {
	return DocumentPath(slug) + "/content"
}
