package server

import (
	"net/http"

	"github.com/containerflow/containerflow/pkg/oauth2"
	"github.com/gin-gonic/gin"
)

// GithubOauthRedirect redirect to github oauth login page
func (s *Server) GithubOauthRedirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, oauth2.RedirectToGithubOauthLoginURL(s.GithubClientID, s.GithubRedirectURI))
}

//GithubOAuthCallback OAuth Callback
func (s *Server) GithubOAuthCallback(c *gin.Context) {
	resp, err := oauth2.GithubOauthCallback(s.GithubClientID, s.GithubClientSecret, c.Query("code"), s.GithubRedirectURI, c.Query("state"))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             resp.Error,
			"error_description": resp.ErrorDescription,
			"error_uri":         resp.ErrorURI,
		})
		return
	}

	user, err := oauth2.FetchGithubUser(resp.AccessToken)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": resp.AccessToken,
		"tokenType":   resp.TokenType,
		"scope":       resp.Scope,
		"userType":    "github",
		"name":        user.Name,
		"email":       user.Email,
		"avatarURL":   user.AvatarURL,
		"login":       user.Login,
	})
}

//Health system health check endpoints
func (s *Server) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
