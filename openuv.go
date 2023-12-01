package openuv

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
	UV               float64          `json:"uv"`
	UVTime           string           `json:"uv_time"`
	UVMax            float64          `json:"uv_max"`
	UVMaxTime        string           `json:"uv_max_time"`
	Ozone            float64          `json:"ozone"`
	OzoneTime        string           `json:"ozone_time"`

	SafeExposureTime SafeExposureTime `json:"safe_exposure_time"`
	SunInfo          SunInfo          `json:"sun_info"`
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
	SunTimes    SunTimes    `json:"sun_times"`
	SunPosition SunPosition `json:"sunt_position"`
}

/* Struct containing information on important sun positions. */
type SunTimes struct {
	Sunrise       string `json:"sunrise"`
	SunriseEnd    string `json:"sunriseEnd"`
	GoldenHourEnd string `json:"goldenHourEnd"`
	SolarNoon     string `json:"solarNoon"`
	GoldenHour    string `json:"goldenHour"`
	SunsetStart   string `json:"sunsetStart"`
	Sunset        string `json:"sunset"`
	Dusk          string `json:"dusk"`
	NauticalDusk  string `json:"nauticalDusk"`
	Night         string `json:"night"`
	Nadir         string `json:"nadir"`
	NightEnd      string `json:"nightEnd"`
	NauticalDawn  string `json:"nauticalDawn"`
	Dawn          string `json:"dawn"`
}

type SunPosition struct {
	Azimuth  float64 `json:"azimuth"`
	Altitude float64 `json:"altitude"`
}

/* Requests the report from OpenUV API.*/
func GetUVReport(client *http.Client, lat, lon float64, keyOUV string) (*UVReport, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(OPEN_UV_URL, lat, lon), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-access-token", keyOUV)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var uvReport UVReport

	err = json.NewDecoder(resp.Body).Decode(&uvReport)

	return &uvReport, err
}
