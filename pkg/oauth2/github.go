package oauth2

import (
	"fmt"

	"github.com/containerflow/containerflow/pkg/rest"
)

//GithubOauthResp Github Oauth2 Response
type GithubOauthResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
}

// GithubUser Github User
type GithubUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Email     string `json:"email"`
}

// FetchGithubUser fetch Github User Info
func FetchGithubUser(accessToken string) (user GithubUser, err error) {
	apiURL := fmt.Sprintf("https://api.github.com/user?access_token=%s", accessToken)
	_, err = rest.Get(apiURL, rest.RequestConfig{}, &user)
	return
}

// RedirectToGithubOauthLoginURL Get Github Oauth2 Login Page Redirect URL
func RedirectToGithubOauthLoginURL(clientID string, redirectURI string) (redirect string) {
	redirect = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=repo,user&state=state", clientID, redirectURI)
	return
}

// GithubOauthCallback callback after github user accept login
func GithubOauthCallback(clientID string, clientSecret string, code string, redirectURI string, state string) (oauthResp GithubOauthResp, err error) {
	var postURL = fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s", clientID, clientSecret, code, redirectURI, state)
	_, err = rest.Post(postURL, rest.RequestConfig{}, &oauthResp)
	if err != nil {
		return
	}
	return
}
