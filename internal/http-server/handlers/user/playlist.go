package user

import (
	resp "SpotifySorter/internal/api/response"
	"SpotifySorter/internal/lib/client/spotify"
	sl "SpotifySorter/internal/lib/logger/slog"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

func GetAllPlaylists(log *slog.Logger, user User) http.HandlerFunc {
	type Response struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			Collaborative bool   `json:"collaborative"`
			Description   string `json:"description"`
			ExternalUrls  struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			Id     string `json:"id"`
			Images []struct {
				Url    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name  string `json:"name"`
			Owner struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Followers struct {
					Href  string `json:"href"`
					Total int    `json:"total"`
				} `json:"followers"`
				Href        string `json:"href"`
				Id          string `json:"id"`
				Type        string `json:"type"`
				Uri         string `json:"uri"`
				DisplayName string `json:"display_name"`
			} `json:"owner"`
			Public     bool   `json:"public"`
			SnapshotId string `json:"snapshot_id"`
			Tracks     struct {
				Href  string `json:"href"`
				Total int    `json:"total"`
			} `json:"tracks"`
			Type string `json:"type"`
			Uri  string `json:"uri"`
		} `json:"items"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Received request on /user/playlist")
		const op = "handlers.user.GetAllPlaylists"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		email := r.URL.Query().Get("email")
		if email == "" {
			log.Error("email parameter is missing")
			render.JSON(w, r, resp.Error("email parameter is required"))
			return
		}

		log.Info("request decoded", slog.String("email", email))

		userData, err := user.GetUser(email)
		if err != nil {
			log.Error("failed to get user data", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get user data"))
			return
		}

		response, err := spotify.GetRequest(log, userData.AccessToken, "users/"+userData.IdSpotify+"/playlists")
		if err != nil {
			log.Error("failed to get all playlists from Spotify", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get all playlists from Spotify"))
			return
		}

		var playlists Response
		if err := json.Unmarshal(response, &playlists); err != nil {
			log.Error("failed to decode response", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode response"))
			return
		}

		render.JSON(w, r, playlists)
	}
}
