package server

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
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

	session, err := mgo.Dial("localhost:27017/test")

	if err != nil {
		log.Fatalln("mongodb connection error", err)
	}

	c := session.DB("").C("sessions")
	store := sessions.NewMongoStore(c, 3600, true, []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// swagger:route GET /ping healthcheck endpoint
	//
	// Check the containerFlow API Server is running
	//
	// You can check the system health status
	//
	//     Responses:
	//       200: pong
	r.GET("/ping", s.Health)
	r.GET("/oauth_redirect", s.GithubOauthRedirect)

	// swagger:route GET /oauth_callback oauth endpoint
	//
	// Github Oauth Callback Endpoint
	//
	// you can get the oauth user info
	//
	//     Responses:
	//       default: error
	//       500: error
	//       400: error
	//       200: oauthuser
	r.GET("/oauth_callback", s.GithubOAuthCallback)
	r.Run()
}
