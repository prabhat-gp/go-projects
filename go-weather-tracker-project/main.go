package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfig struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (*apiConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c *apiConfig
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}

func query(city string, apiKey string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weatherData{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func main() {
	apiConfigData, err := loadApiConfig(".apiConfig")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city, apiConfigData.OpenWeatherMapApiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
