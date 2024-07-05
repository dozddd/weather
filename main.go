package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	points := []GeoPoint{
		{
			lat: 55.751244,
			lng: 37.618423,
		},
		{
			lat: 48.864716,
			lng: 2.349014,
		},
		{
			lat: 61.25,
			lng: 73.4166667,
		},
	}

	resps := make(chan *WeatherResp, len(points))

	for _, p := range points {
		go func(p GeoPoint) {
			r, err := p.GetWeather()
			if err != nil {
				log.Fatalf("failed get wether: %v", err)
			}
			resps <- r
		}(p)
	}

	var sumTemp float64
	for i := 0; i < len(points); i++ {
		r := <-resps
		sumTemp += r.Current.Temperature2M
	}

	fmt.Println(sumTemp / float64(len(points)))
}

type WeatherResp struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationTimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentUnits         struct {
		Time          string `json:"time"`
		Interval      string `json:"interval"`
		Temperature2M string `json:"temperature_2m"`
	} `json:"current_units"`
	Current struct {
		Time          string  `json:"time"`
		Interval      int     `json:"interval"`
		Temperature2M float64 `json:"temperature_2m"`
	} `json:"current"`
}

type GeoPoint struct {
	lng float64
	lat float64
}

func (p GeoPoint) GetWeather() (*WeatherResp, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
		p.lat,
		p.lng,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	r := WeatherResp{}
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &r, nil
}
