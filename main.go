package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type apiConfidData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfidData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return apiConfidData{}, err
	}
	var c apiConfidData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfidData{}, err
	}
	return c, nil
}
//middleware
func enableCORS(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
func hello(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	w.Write([]byte("hello from go!\n"))
}
func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	encodedCity := url.QueryEscape(city)
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + encodedCity)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return weatherData{}, fmt.Errorf("API request failed with status:%s", resp.Status)
	}
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil

}
func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			enableCORS(w, r)
				if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})
	http.ListenAndServe(":8083", nil)
}
