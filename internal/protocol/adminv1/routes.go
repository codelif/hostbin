package adminv1

const (
	BasePath               = "/api/v1"
	HealthPath             = BasePath + "/health"
	DocumentsCollection    = BasePath + "/documents"
	DocumentPathPattern    = "/documents/:slug"
	DocumentContentPattern = "/documents/:slug/content"
)
