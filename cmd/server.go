package main

import "github.com/tuananhlai/brevity-go/internal/config"

func runServer() {
	cfg := config.MustLoadConfig()
	println(cfg.Database.URL)
	println("Server stopped.")
}
