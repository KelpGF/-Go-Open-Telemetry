package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/KelpGF/Go-Observability/internal/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Body struct {
	ZipCode string `json:"cep"`
}

func (b Body) isValid() bool {
	return len(b.ZipCode) == 8
}

func Validate(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	var body Body

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		errorMessage := newResponseError("Invalid body")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
	}

	if !body.isValid() {
		errorMessage := newResponseError("Invalid ZipCode")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	response, err := services.HttpRequest(
		ctx,
		"http://"+os.Getenv("API_DNS")+"/zip-code/weather?zipcode="+body.ZipCode,
	)

	if err != nil {
		errorMessage := newResponseError("Error on request ZipCode data:" + err.Error())

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	output := WeatherResult{}
	err = json.Unmarshal(response.Data, &output)

	if err != nil {
		errorMessage := newResponseError("Error on request ZipCode data")

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
