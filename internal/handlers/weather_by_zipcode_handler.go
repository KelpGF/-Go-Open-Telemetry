package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KelpGF/Go-Observability/internal/services"
)

type SearchCEPResult struct {
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

type WeatherResult struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func WeatherByCepHandler(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("zipcode")
	w.Header().Set("Content-Type", "application/json")

	zipCodeData, err := services.GetZipCodeData(zipCode)

	if err != nil {
		errorString := err.Error()
		statusCode := http.StatusInternalServerError

		if err.Error() == "invalid zipcode" {
			statusCode = http.StatusUnprocessableEntity
		}

		if err.Error() == "can not find zipcode" {
			statusCode = http.StatusNotFound
		}

		errorMessage := newResponseError(errorString)

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	weatherData, err := services.GetWeatherData(zipCodeData.Localidade)

	if err != nil {
		errorMessage := newResponseError("Error on weather request")

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	result := WeatherResult{
		City:  zipCodeData.Localidade,
		TempC: weatherData.Current.TempC,
		TempF: weatherData.Current.TempF,
		TempK: weatherData.Current.TempC + 273.15,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
