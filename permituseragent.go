// Package traefik_plugin_permituseragent a plugin to permit traffic based on a User-Agent string.
package traefik_plugin_permituseragent

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

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
}

// New creates and returns a plugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	fmt.Println("Creating middleware")

	if config.UserAgent == "" {
		return nil, fmt.Errorf(name, " must have UserAgent value set.")
	}

	fmt.Printf("Loading %s plugin with settings: url -> %s useragent -> %s", name, config.Url, config.UserAgent)

	return &permitUserAgent{
		name:    name,
		next:    next,
		userAgent: config.UserAgent,
		url: config.Url,
	}, nil
}

func (p *permitUserAgent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req != nil {
		fmt.Printf("Checking UserAgent for connection from: %s", req.RemoteAddr)
		userAgent := req.UserAgent()

		// if the request useragent doesn't match our defined useragent
		// then redirect to the defined url
		if userAgent != p.userAgent {
			log.Printf("Redirecting User-Agent: '%s' to URL: '%s'", userAgent, p.url)
			rw.Header().Set("Location", p.url)
			rw.WriteHeader(http.StatusFound)
			return
		} else {
			fmt.Printf("Successful UserAgent match: %s", userAgent)
		}
	}

	p.next.ServeHTTP(rw, req)
}
