package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
)

// Super dumb HTTP server, something for the `observer` to start and control.
func main() {

	var (
		flagAddress string
		flagName    string
	)

	pflag.StringVarP(&flagAddress, "address", "a", ":8080", "address for the server to use")
	pflag.StringVarP(&flagName, "name", "n", "default-rubbish-name", "name for the server to use")

	pflag.Parse()

	srv := echo.New()

	srv.GET("/version", versionHandler)
	srv.GET("/name", nameHandler(flagName))

	err := srv.Start(flagAddress)
	if err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}

func versionHandler(ctx echo.Context) error {
	log.Printf("received version request")
	return ctx.String(http.StatusOK, "1.0.0")
}

func nameHandler(name string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Printf("received name request")
		return ctx.String(http.StatusOK, name)
	}
}
