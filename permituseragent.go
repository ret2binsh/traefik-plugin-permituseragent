// Package traefik_plugin_permituseragent a plugin to permit traffic based on a User-Agent string.
package traefik_plugin_permituseragent

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/traefik/traefik/v3/pkg/middlewares"
)

const typeName = "PermitUserAgent"

// Config holds the plugin configuration.
type Config struct {
	UserAgent string `json:"userAgent,omitempty"`
	Url string `json:"url,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{
		UserAgent: "",
		Url: "https://google.com",
	}
}

type permitUserAgent struct {
	name    string
	next    http.Handler
	userAgent string
	url string
	logger *zerolog.Logger
}

// New creates and returns a plugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	logger := middlewares.GetLogger(ctx, name, typeName)
	logger.Debug().Msg("Creating middleware")

	if config.UserAgent == "" {
		return nil, fmt.Errorf(name, " must have UserAgent value set.")
	}

	logger.Debug().Msgf("Loading %s plugin with settings: url -> %s useragent -> %s", name, config.Url, config.UserAgent)

	return &permitUserAgent{
		name:    name,
		next:    next,
		userAgent: config.UserAgent,
		url: config.Url,
		logger: logger,
	}, nil
}

func (p *permitUserAgent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req != nil {
		p.logger.Debug().Msgf("Checking UserAgent for connection from: %s", req.RemoteAddr)
		userAgent := req.UserAgent()

		// if the request useragent doesn't match our defined useragent
		// then redirect to the defined url
		if userAgent != p.userAgent {
			p.logger.Warn().Msgf("Redirecting User-Agent: '%s' to URL: '%s'", userAgent, p.url)
			rw.Header().Set("Location", p.url)
			rw.WriteHeader(http.StatusFound)
			return
		} else {
			p.logger.Debug().Msgf("Successful UserAgent match: %s", userAgent)
		}
	}

	p.next.ServeHTTP(rw, req)
}
