package user

import (
	resp "SpotifySorter/internal/api/response"
	"SpotifySorter/internal/lib/client/spotify"
	sl "SpotifySorter/internal/lib/logger/slog"
	UserModel "SpotifySorter/models"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type User interface {
	SaveUser(email, accessToken, country, name, href, idSpotify, product, uri string) (*UserModel.User, error)
	GetUserByEmail(email string) (*UserModel.User, error)
	UpdateUser(email, accessToken, country, name, href, product, uri string) (*UserModel.User, error)
}

func AuthUser(log *slog.Logger, user User) http.HandlerFunc {
	type Request struct {
		Code  string `json:"code" validate:"required"`
		State string `json:"state"`
	}
	type Response struct {
		resp.Response
	}
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.AuthUser"
		log = log.With(slog.String("op", op))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		accessCredentials, err := sendCode(log, req.Code)
		if err != nil {
			log.Error("failed to send code user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to send code user"))
			return
		}

		userData, err := getUserData(log, accessCredentials)
		if err != nil {
			log.Error("failed to get user data", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get user data"))
			return
		}
		_, err = user.GetUserByEmail(userData.Email)
		if errors.Is(err, sql.ErrNoRows) {
			savedUser, err := user.SaveUser(userData.Email, accessCredentials.AccessToken, userData.Country, userData.Name, userData.Href, userData.IdSpotify, userData.Product, userData.Uri)
			if err != nil {
				log.Error("failed to save user", sl.Err(err))
				render.JSON(w, r, resp.Error("failed to save user"))
				return
			}
			userData.Id = savedUser.Id

			render.JSON(w, r, userData)
			return
		} else {
			updatedUser, err := user.UpdateUser(userData.Email, accessCredentials.AccessToken, userData.Country, userData.Name, userData.Href, userData.Product, userData.Uri)
			if err != nil {
				log.Error("failed to update user", sl.Err(err))
				render.JSON(w, r, resp.Error("failed to update user"))
				return
			}
			render.JSON(w, r, updatedUser)
			return
		}
	}
}

func sendCode(log *slog.Logger, code string) (*UserModel.AccessCredentials, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))
	data.Set("grant_type", "authorization_code")

	reqToSpotify, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("failed to create request", sl.Err(err))
		return nil, errors.New("failed to create request")
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Error("SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET is not set")
		return nil, errors.New("SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET is not set")
	}
	auth := clientID + ":" + clientSecret
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	reqToSpotify.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqToSpotify.Header.Set("Authorization", "Basic "+encodedAuth)

	client := &http.Client{}
	resp, err := client.Do(reqToSpotify)
	if err != nil {
		log.Error("failed to send request to Spotify", sl.Err(err))
		return nil, errors.New("failed to send request to Spotify")
	}
	defer resp.Body.Close()

	var userData UserModel.AccessCredentials
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		log.Error("failed to decode response from Spotify", sl.Err(err))
		return nil, errors.New("failed to decode response from Spotify")
	}

	return &userData, nil
}

func getUserData(log *slog.Logger, accessCredentials *UserModel.AccessCredentials) (*UserModel.User, error) {
	response, err := spotify.GetRequest(log, accessCredentials.AccessToken, "me")
	if err != nil {
		log.Error("failed to get response from Spotify", sl.Err(err))
		return nil, errors.New("failed to get response from Spotify")
	}

	var user UserModel.User
	if err := json.Unmarshal(response, &user); err != nil {
		log.Error("failed to decode response", sl.Err(err))
		return nil, errors.New("failed to decode response from Spotify")
	}

	return &user, nil
}
