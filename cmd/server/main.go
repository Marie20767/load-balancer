package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func run(port string) error {
	e := echo.New()
	e.GET("/hello/:name", func(c echo.Context) error {
		name := c.Param("name")

		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Hello, %s!", name),
		})
	})

	return e.Start(fmt.Sprintf(":%s", port))
}

func main() {
	port := flag.String("port", "", "Port to listen on")
	flag.Parse()

	if *port == "" {
		log.Fatal("Port is required.")
	}

	if err := run(*port); err != nil {
		log.Fatal("server closed: ", err)
	}
}
