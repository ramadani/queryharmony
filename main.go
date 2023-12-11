package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// pingHandler handles the "/ping" endpoint and returns a "Pong!" response.
func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Pong!")
}

func main() {
	// Declare and set the default value for the port flag
	port := flag.Int("port", 8080, "Port number")
	flag.Parse()

	// Create an instance of Echo
	e := echo.New()

	// Set the handler for the "/ping" endpoint
	e.GET("/ping", pingHandler)

	// Run the server
	fmt.Printf("Server is listening on :%d...\n", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
