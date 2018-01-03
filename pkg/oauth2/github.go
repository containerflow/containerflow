package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// RedirectToGithubOauthLoginURL Get Github Oauth2 Login Page Redirect URL
func RedirectToGithubOauthLoginURL(clientID string, redirectURI string) (redirect string) {
	redirect = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=repo,user&state=state", clientID, redirectURI)
	return
}

// GithubOauthCallback callback after github user accept login
func GithubOauthCallback(clientId string, clientSecret string, code string, redirectURI string, state string) (oauthResp GithubOauthResp, err error) {
	var postURL = fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s", clientId, clientSecret, code, redirectURI, state)

	client := &http.Client{}
	req, err := http.NewRequest("POST", postURL, nil)

	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	json.Unmarshal(body, &oauthResp)
	return
}
