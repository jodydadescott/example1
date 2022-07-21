package operator

import "net/http"

// Config ...
type Config struct {
	PrismaAPI       string
	PrismaLabel     string
	PrismaNamespace string
	LabelSelectors  []string
	HTTPClient      *http.Client
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}

// SetPrismaAPI ...
func (t *Config) SetPrismaAPI(prismaAPI string) *Config {
	t.PrismaAPI = prismaAPI
	return t
}

// SetPrismaLabel ...
func (t *Config) SetPrismaLabel(prismaLabel string) *Config {
	t.PrismaLabel = prismaLabel
	return t
}

// SetPrismaNamespace ...
func (t *Config) SetPrismaNamespace(prismaNamespace string) *Config {
	t.PrismaNamespace = prismaNamespace
	return t
}

// SetHTTPClient ...
func (t *Config) SetHTTPClient(httpClient *http.Client) *Config {
	t.HTTPClient = httpClient
	return t
}

// AddLabelSelectors ...
func (t *Config) AddLabelSelectors(labelSelectors ...string) *Config {
	t.LabelSelectors = append(t.LabelSelectors, labelSelectors...)
	return t
}
