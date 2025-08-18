package main

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func loadBalancer(c echo.Context) error {
	// TODO: route to one of the servers with round robin or first without
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	port := os.Getenv("PORT")
	if port == "" {
		return errors.New("not all environment variables are set")
	}

	e := echo.New()
	e.GET("*", loadBalancer)

	return e.Start(":" + port)
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}
