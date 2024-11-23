package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/KelpGF/Go-Observability/configs"
	"github.com/KelpGF/Go-Observability/internal/handlers"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	println("Starting server...")

	appName := os.Getenv("APP_NAME")
	port := os.Getenv("PORT")

	if err := run(appName, port); err != nil {
		log.Fatalln(err)
	}
}

func run(appName string, port string) (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	println("Setting up OpenTelemetry SDK...")
	otelShutdown, err := configs.SetupOTelSDK(ctx, appName)
	if err != nil {
		println("Failed to setup OpenTelemetry SDK.")
		return
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	srv := &http.Server{
		Addr:         ":" + port,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		println("Starting server...")
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	println("Shutting down server...")
	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	handleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	})

	handleFunc("/zip-code/validate", handlers.Validate)
	handleFunc("/zip-code/weather", handlers.WeatherByCepHandler)

	handler := otelhttp.NewHandler(mux, "/")
	return handler
}
