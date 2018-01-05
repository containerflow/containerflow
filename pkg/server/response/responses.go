package response

//Pong System Pong Response
//swagger:response pong
type Pong struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

//Error system error response
//swagger:response error
type Error struct {
	Error            string `json:"message"`
	ErrorDescription string `json:"error_description"`
}

// OauthUser github oauth user response
//swagger:response oauthuser
type OauthUser struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	Scope       string `json:"scope"`
	UserType    string `json:"userType"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatarURL"`
	Login       string `json:"login"`
}
