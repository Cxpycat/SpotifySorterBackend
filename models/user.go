package user

type User struct {
	Id           int64
	Country      string `json:"country"`
	Name         string `json:"display_name"`
	AccessToken  string `json:"access_token,omitempty"`
	Email        string `json:"email"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href      string `json:"href"`
	IdSpotify string `json:"id"`
	Product   string `json:"product"`
	Uri       string `json:"uri"`
}

type AccessCredentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
