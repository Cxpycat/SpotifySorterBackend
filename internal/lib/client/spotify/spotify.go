package spotify

import (
	sl "SpotifySorter/internal/lib/logger/slog"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

func GetRequest(log *slog.Logger, accessToken, endpoint string) ([]byte, error) {
	reqToSpotify, err := http.NewRequest("GET", "https://api.spotify.com/v1/"+endpoint, nil)
	if err != nil {
		log.Error("failed to create request", sl.Err(err))
		return nil, errors.New("failed to create request")
	}

	reqToSpotify.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(reqToSpotify)
	if err != nil {
		log.Error("failed to send request", slog.String("url", "https://api.spotify.com/v1/"+endpoint), sl.Err(err))
		return nil, errors.New("failed to send request to Spotify")
	}
	defer resp.Body.Close()

	log.Info("response received", slog.String("url", "https://api.spotify.com/v1/"+endpoint), slog.Int("status", resp.StatusCode))

	if resp.StatusCode == http.StatusForbidden {
		log.Error("unauthorized")
		return nil, errors.New("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("bad request", slog.Any("body", resp.Body))
		return nil, errors.New("bad request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body", sl.Err(err))
		return nil, errors.New("failed to read response body")
	}

	return body, nil
}
