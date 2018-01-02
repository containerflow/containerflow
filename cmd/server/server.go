package main

import (
	"flag"

	"github.com/containerflow/containerflow/pkg/server"
	"github.com/gobike/envflag"
)

var (
	githubClientID     string
	githubClientSecret string
	githubRedirectURI  string
)

func init() {
	flag.StringVar(&githubClientID, "githubclientid", "clientid", "Github OAuth App Client Id")
	flag.StringVar(&githubClientSecret, "githubclientsecret", "clientsecret", "Github OAuth App Client Secret")
	flag.StringVar(&githubRedirectURI, "githubredirecturi", "http://127.0.0.1:8080/oauth_callback", "Github OAuth Client Redirect URL")
}

func main() {
	envflag.Parse()
	server := server.Server{
		GithubClientID:     githubClientID,
		GithubClientSecret: githubClientSecret,
		GithubRedirectURI:  githubRedirectURI,
	}
	server.Start()
}
