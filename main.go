package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	//"strconv"
	"strings"
)

const baseURL = "https://api.open-meteo.com/v1/forecast"
const geonamesUsername = "jrh230"

type WeatherData struct {
	Daily                DailyData  `json:"daily"`
	DailyUnits           DailyUnits `json:"daily_units"`
	Elevation            float64    `json:"elevation"`
	GenerationtimeMs     float64    `json:"generationtime_ms"`
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	Timezone             string     `json:"timezone"`
	TimezoneAbbreviation string     `json:"timezone_abbreviation"`
	UtcOffsetSeconds     int        `json:"utc_offset_seconds"`
}

type DailyData struct {
	Temperature2mMax []float64 `json:"temperature_2m_max"`
	Temperature2mMin []float64 `json:"temperature_2m_min"`
	Time             []string  `json:"time"`
}

type DailyUnits struct {
	Temperature2mMax string `json:"temperature_2m_max"`
	Temperature2mMin string `json:"temperature_2m_min"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: go_weather_cli <latitude,longitude>")
		os.Exit(1)
	}

	location := args[1]
	weatherData, err := fetchWeatherData(location)
	if err != nil {
		fmt.Printf("Failed to fetch weather data: %v", err)
		os.Exit(1)
	}

	printWeatherData(location, weatherData)
}

func fetchWeatherData(location string) (*WeatherData, error) {
	coords := strings.Split(location, ",")
	latitude, longitude := coords[0], coords[1]

	timezoneAPIURL := fmt.Sprintf(
		"http://api.geonames.org/timezoneJSON?lat=%s&lng=%s&username=%s",
		latitude, longitude, geonamesUsername)

	timezoneData, err := getJSON(timezoneAPIURL)
	if err != nil {
		return nil, err
	}

	timezoneId, ok := timezoneData["timezoneId"].(string)
	if !ok {
		fmt.Println("Failed to fetch timezone data")
		os.Exit(1)
	}

	url := fmt.Sprintf(
		"%s?latitude=%s&longitude=%s&timezone=%s&daily=temperature_2m_min,temperature_2m_max",
		baseURL, latitude, longitude, timezoneId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, err
	}

	return &weatherData, nil
}

func printWeatherData(location string, weatherData *WeatherData) {
	fmt.Printf("Weather data for %s:\n", location)
	fmt.Printf("Timezone: %s\n", weatherData.Timezone)
	fmt.Printf("Elevation: %f meters\n", weatherData.Elevation)
	fmt.Printf("Generation Time (ms): %f\n", weatherData.GenerationtimeMs)
	fmt.Printf("Latitude: %f\n", weatherData.Latitude)
	fmt.Printf("Longitude: %f\n", weatherData.Longitude)
	fmt.Printf("Timezone Abbreviation: %s\n", weatherData.TimezoneAbbreviation)
	fmt.Printf("UTC Offset (seconds): %d\n", weatherData.UtcOffsetSeconds)

	for i, time := range weatherData.Daily.Time {
		fmt.Printf("Date: %s\n", time)
		maxTempF := weatherData.Daily.Temperature2mMax[i]*9.0/5.0 + 32.0
		minTempF := weatherData.Daily.Temperature2mMin[i]*9.0/5.0 + 32.0
		fmt.Printf("Max Temperature: %f F\n", maxTempF)
		fmt.Printf("Min Temperature: %f F\n\n", minTempF)
	}
}

func getJSON(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
