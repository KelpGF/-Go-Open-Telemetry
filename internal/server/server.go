package server

import (
	"crypto/tls"
	"net/http"

	"github.com/KelpGF/Go-Observability/internal/handlers"
)

type ServerHttp struct{}

func (*ServerHttp) Run() {
	println("Starts server on port 8080")

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Kelp Weather By ZipeCode!"))
	})

	http.HandleFunc("/zip-code/validate", handlers.Validate)
	http.HandleFunc("/zip-code/weather", handlers.WeatherByCepHandler)

	http.ListenAndServe(":8080", nil)
}
