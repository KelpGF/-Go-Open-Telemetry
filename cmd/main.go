package main

import (
	"github.com/KelpGF/Go-Observability/internal/server"
)

func main() {
	server := server.ServerHttp{}

	server.Run()
}
