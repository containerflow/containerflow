package server

import (
	"github.com/gin-gonic/gin"
)

// Server container web server
type Server struct {
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURI  string
}

// Start start web server
func (s *Server) Start() {
	r := gin.Default()
	r.GET("/ping", s.Health)
	r.GET("/oauth_redirect", s.GithubOauthRedirect)
	r.GET("/oauth_callback", s.GithubOAuthCallback)
	r.Run()
}
