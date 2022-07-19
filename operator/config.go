package operator

import "net/http"

// Config ...
type Config struct {
	PrismaAPI       string
	PrismaLabel     string
	PrismaNamespace string
	HTTPClient      *http.Client
}
