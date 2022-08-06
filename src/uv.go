package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A template URL for requesting data from OpenUV API
const OPEN_UV_URL string = "https://api.openuv.io/api/v1/uv?lat=%f&lng=%f"

/* The top structure of the OpenUV API response. */
type UVReport struct {
	Result `json:"result"`
}

/* Result returned from OpenUV API call. */
type Result struct {
	UV       float32 `json:"uv"`
	UVMax    float32 `json:"uv_max"`

	Ozone    float32 `json:"ozone"`

	SafeExposureTime `json:"safe_exposure_time"`

	SunInfo          `json:"sun_info"`
}

/* Safe exposure time in minutes for different skin types (Fitzpatrick scale). */
type SafeExposureTime struct {
	ST1 int `json:"st1"` // very fair skin, white
	ST2 int `json:"st2"` // fair skin, white
	ST3 int `json:"st3"` // fair skin, cream white
	ST4 int `json:"st4"` // olive skin
	ST5 int `json:"st5"` // brown skin
	ST6 int `json:"st6"` // black skin
}

/* The top structure for the data related to the Sun. */
type SunInfo struct {
	SunTimes `json:"sun_times"`
}

/* Struct containing information on important sun positions. */
type SunTimes struct {
	Sunrise    string `json:"sunrise"`
	SolarNoon  string `json:"solarNoon"`
	Sunset     string `json:"sunset"`
	Night      string `json:"night"`

	GoldenHour string `json:"goldenHour"`
	GHMorning  string `json:"goldenHourEnd"`
}

/* Formats the UVReport struct into string. */
func (uv *UVReport) ToString() string {
	return fmt.Sprintf(
		"UV Index:\n"                 + 
		"  Current: %6.2f\n"          + 
		"  Max:     %6.2f\n"          + 
		"  Ozone:   %6.2f\n\n"        + 
		"Sunrise: %32s\n"             + 
		"Solar Noon: %29s\n"          + 
		"Sunset: %33s\n"              + 
		"Night: %34s\n"               + 
		"Golden Hour: %28s\n"         + 
		"Morning GH ends: %s\n\n"     + 
		"Safe Exposure Time [min]:\n" + 
		"  1: %5d   |   4: %5d\n"     + 
		"  2: %5d   |   5: %5d\n"     + 
		"  3: %5d   |   6: %5d", 
		uv.UV, uv.UVMax,
		uv.Ozone,
		uv.Sunrise,
		uv.SolarNoon,
		uv.Sunset,
		uv.Night, 
		uv.GoldenHour,
		uv.GHMorning,
		uv.ST1, uv.ST4, 
		uv.ST2, uv.ST5, 
		uv.ST3, uv.ST6,
	)
}

/* Requests the report from OpenUV API.*/
func GetUVReport(client *http.Client, loc *Location, s *Settings) (*UVReport, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(OPEN_UV_URL, loc.Lat, loc.Lon), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-access-token", s.OpenUVKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var uvReport UVReport

	err = json.NewDecoder(resp.Body).Decode(&uvReport)
	
	return &uvReport, err
}