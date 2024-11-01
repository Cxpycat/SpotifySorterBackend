package user

import (
	resp "SpotifySorter/internal/api/response"
	"SpotifySorter/internal/lib/client/spotify"
	sl "SpotifySorter/internal/lib/logger/slog"
	"SpotifySorter/internal/storage"
	"encoding/json"
	"github.com/go-chi/chi/v5"
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
		const op = "handlers.playlist.GetAllPlaylists"

		log = log.With(slog.String("op", op))

		email := r.URL.Query().Get("email")
		if email == "" {
			log.Error("email parameter is missing")
			render.JSON(w, r, resp.Error("email parameter is required"))
			return
		}

		userData, err := user.GetUserByEmail(email)
		if err != nil {
			log.Error(storage.ErrUserNotFound.Error(), sl.Err(err))
			render.JSON(w, r, resp.Error(storage.ErrUserNotFound.Error()))
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

func GetPlaylistById(log *slog.Logger, user User) http.HandlerFunc {
	type Response struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			AddedAt string `json:"added_at"`
			AddedBy struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Followers struct {
					Href  string `json:"href"`
					Total int    `json:"total"`
				} `json:"followers"`
				Href string `json:"href"`
				Id   string `json:"id"`
				Type string `json:"type"`
				Uri  string `json:"uri"`
			} `json:"added_by"`
			IsLocal bool `json:"is_local"`
			Track   struct {
				Album struct {
					AlbumType        string   `json:"album_type"`
					TotalTracks      int      `json:"total_tracks"`
					AvailableMarkets []string `json:"available_markets"`
					ExternalUrls     struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href   string `json:"href"`
					Id     string `json:"id"`
					Images []struct {
						Url    string `json:"url"`
						Height int    `json:"height"`
						Width  int    `json:"width"`
					} `json:"images"`
					Name                 string `json:"name"`
					ReleaseDate          string `json:"release_date"`
					ReleaseDatePrecision string `json:"release_date_precision"`
					Restrictions         struct {
						Reason string `json:"reason"`
					} `json:"restrictions"`
					Type    string `json:"type"`
					Uri     string `json:"uri"`
					Artists []struct {
						ExternalUrls struct {
							Spotify string `json:"spotify"`
						} `json:"external_urls"`
						Href string `json:"href"`
						Id   string `json:"id"`
						Name string `json:"name"`
						Type string `json:"type"`
						Uri  string `json:"uri"`
					} `json:"artists"`
				} `json:"album"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					Id   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					Uri  string `json:"uri"`
				} `json:"artists"`
				AvailableMarkets []string `json:"available_markets"`
				DiscNumber       int      `json:"disc_number"`
				DurationMs       int      `json:"duration_ms"`
				Explicit         bool     `json:"explicit"`
				ExternalIds      struct {
					Isrc string `json:"isrc"`
					Ean  string `json:"ean"`
					Upc  string `json:"upc"`
				} `json:"external_ids"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href       string `json:"href"`
				Id         string `json:"id"`
				IsPlayable bool   `json:"is_playable"`
				LinkedFrom struct {
				} `json:"linked_from"`
				Restrictions struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Name        string `json:"name"`
				Popularity  int    `json:"popularity"`
				PreviewUrl  string `json:"preview_url"`
				TrackNumber int    `json:"track_number"`
				Type        string `json:"type"`
				Uri         string `json:"uri"`
				IsLocal     bool   `json:"is_local"`
			} `json:"track"`
		} `json:"items"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.playlist.GetPlaylistById"

		log = log.With(slog.String("op", op))

		email := r.URL.Query().Get("email")
		if email == "" {
			log.Error("email parameter is missing")
			render.JSON(w, r, resp.Error("email parameter is required"))
			return
		}

		userData, err := user.GetUserByEmail(email)
		if err != nil {
			log.Error(storage.ErrUserNotFound.Error(), sl.Err(err))
			render.JSON(w, r, resp.Error(storage.ErrUserNotFound.Error()))
			return
		}

		id := chi.URLParam(r, "id")

		response, err := spotify.GetRequest(log, userData.AccessToken, "playlists/"+id+"/tracks/")
		if err != nil {
			log.Error("Error getting playlist by ID", sl.Err(err))
			http.Error(w, "Error getting playlist by ID", http.StatusInternalServerError)
			return
		}

		var playlist Response
		if err := json.Unmarshal(response, &playlist); err != nil {
			log.Error("failed to decode response", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode response"))
			return
		}

		render.JSON(w, r, playlist)
	}
}
