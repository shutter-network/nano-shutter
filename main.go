package main

import (
	"context"
	"nano-shutter/dkg"
	"nano-shutter/router"
	"nano-shutter/service"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	pubdkg := dkg.StartDkg()

	srv := service.NewService(pubdkg)
	app := router.NewRouter(context.Background(), srv)
	app.Run("0.0.0.0:" + port)
}
