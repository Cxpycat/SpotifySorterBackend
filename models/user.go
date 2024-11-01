package user

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Id                 int64
	Country            string `json:"country"`
	Name               string `json:"display_name"`
	AccessToken        string `json:"access_token,omitempty"`
	SpotifyAccessToken string `json:"spotify_access_token,omitempty"`
	Email              string `json:"email"`
	IdSpotify          string `json:"id"`
	Product            string `json:"product"`
}

type AccessTokensByCode struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	Name        string `json:"name"`
	AccessToken string `json:"access_token,omitempty"`
	Email       string `json:"email"`
	IdSpotify   string `json:"id"`
}

type Claims struct {
	jwt.RegisteredClaims
}
