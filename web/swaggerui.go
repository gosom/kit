package web

import (
	"embed"
	"net/http"

	"github.com/ismurov/swaggerui"
)

type SwaggerUIConfig struct {
	SpecName string
	SpecFile string
	SpecFS   embed.FS
	Path     string
}

func NewSwaggerUI(cfg *SwaggerUIConfig) (http.Handler, error) {
	h, err := swaggerui.New(
		[]swaggerui.SpecFile{{
			Name: cfg.SpecName,
			Path: cfg.SpecFile,
		}},
		cfg.SpecFS,
	)
	return h, err
}
