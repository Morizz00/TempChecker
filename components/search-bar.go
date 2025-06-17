package components

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Location struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	State   string  `json:"state,omitempty"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type SearchBar struct {
	APIKey string
	Client *http.Client
}

func NewSearchBar(apiKey string) *SearchBar {
	return &SearchBar{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SearchBar) Search(query string) ([]Location, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, fmt.Errorf("query too short")
	}

	baseURL := "http://api.openweathermap.org/geo/1.0/direct"
	params := url.Values{}
	params.Add("q", query)
	params.Add("limit", "5")
	params.Add("appid", s.APIKey)

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := s.Client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp []struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		State   string  `json:"state,omitempty"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	locations := make([]Location, len(apiResp))
	for i, item := range apiResp {
		locations[i] = Location{
			Name:    item.Name,
			Country: item.Country,
			State:   item.State,
			Lat:     item.Lat,
			Lon:     item.Lon,
		}
	}

	return locations, nil
}

func (s *SearchBar) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	locations, err := s.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(locations)
}
