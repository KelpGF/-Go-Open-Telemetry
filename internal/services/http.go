package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type searchCEPResult struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Error       string `json:"erro"`
}

type searchWeatherResult struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree float64 `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   float64 `json:"humidity"`
		Cloud      float64 `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

type httpRequestOutput struct {
	Code int
	Data []byte
}

func HttpRequest(url string) (httpRequestOutput, error) {
	output := httpRequestOutput{}

	req, err := http.Get(url)
	if err != nil {
		return output, err
	}

	output.Code = req.StatusCode

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return output, err
	}

	output.Data = res

	return output, nil
}

func GetZipCodeData(zipCode string) (searchCEPResult, error) {
	output := searchCEPResult{}

	urlRequest := "https://viacep.com.br/ws/" + zipCode + "/json/"
	result, err := HttpRequest(urlRequest)
	if err != nil {
		return output, err
	}

	statusCode := result.Code
	if statusCode == 400 {
		return output, errors.New("invalid zipcode")
	}

	err = json.Unmarshal(result.Data, &output)
	if err != nil {
		return output, errors.New("error on zipcode data format")
	}

	if output.Error == "true" {
		return output, errors.New("can not find zipcode")
	}

	return output, nil
}

func GetWeatherData(city string) (searchWeatherResult, error) {
	output := searchWeatherResult{}

	urlRequest := "http://api.weatherapi.com/v1/current.json?key=c4a0b6bf6e1342c38f3153503242109&q=" + url.QueryEscape(city)

	result, err := HttpRequest(urlRequest)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(result.Data, &output)
	if err != nil {
		return output, err
	}

	return output, nil
}
