package openuv

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zedran/geoloc"
)

// A template URL for requesting data from OpenUV API
const OPEN_UV_URL string = "https://api.openuv.io/api/v1/uv?lat=%f&lng=%f"

/* The top structure of the OpenUV API response. */
type UVReport struct {
	Result `json:"result"`
}

/* Result returned from OpenUV API call. */
type Result struct {
	UV        float32 `json:"uv"`
	UVMax     float32 `json:"uv_max"`
	Ozone     float32 `json:"ozone"`

	UVTime    string  `json:"uv_time"`
	UVMaxTime string  `json:"uv_max_time"`
	OzoneTime string  `json:"ozone_time"`

	SafeExposureTime  `json:"safe_exposure_time"`

	SunInfo           `json:"sun_info"`
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

/* Requests the report from OpenUV API.*/
func GetUVReport(client *http.Client, loc *geoloc.Location, keyOUV string) (*UVReport, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(OPEN_UV_URL, loc.Lat, loc.Lon), nil)
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
