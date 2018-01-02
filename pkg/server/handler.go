package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GithubOauthRedirect redirect to github oauth login page
func (s *Server) GithubOauthRedirect(c *gin.Context) {
	var redirectURL = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=repo,user&state=state", s.GithubClientID, s.GithubRedirectURI)
	log.Println(redirectURL)
	c.Redirect(http.StatusMovedPermanently, redirectURL)
}

//GithubOAuthCallback OAuth Callback
func (s *Server) GithubOAuthCallback(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "next should post to https://github.com/login/oauth/access_token get access token",
		"code":    c.Query("code"),
	})
}

//Health system health check endpoints
func (s *Server) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
