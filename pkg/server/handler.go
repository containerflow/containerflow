package server

import (
	"net/http"

	"github.com/containerflow/containerflow/pkg/oauth2"
	"github.com/containerflow/containerflow/pkg/server/response"
	"github.com/gin-contrib/sessions"
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
		c.JSON(http.StatusBadRequest, response.Error{
			Error:            resp.Error,
			ErrorDescription: resp.ErrorDescription,
		})
		return
	}

	user, err := oauth2.FetchGithubUser(resp.AccessToken)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, response.Error{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.OauthUser{
		AccessToken: resp.AccessToken,
		TokenType:   resp.TokenType,
		Scope:       resp.Scope,
		UserType:    "github",
		Name:        user.Name,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Login:       user.Login,
	})
}

//Health system health check endpoints
func (s *Server) Health(c *gin.Context) {

	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()

	c.JSON(200, response.Pong{
		Message: "pong",
		Count:   count,
	})
}
