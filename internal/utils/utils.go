package utils

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func CustomErrHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var code int
	var msg string

	if _, isUrlError := err.(*url.Error); isUrlError {
		code = http.StatusBadGateway
		msg = "Bad Gateway: " + err.Error()
	} else {
		code = http.StatusInternalServerError
		msg = err.Error()
	}

	if err := c.JSON(code, map[string]string{"error": msg}); err != nil {
		log.Printf("failed to send JSON response: %v", err)
	}
}

type Config struct {
	Port string
}

func ParseEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, errors.New("port environment variable not set")
	}

	return &Config{
		Port: port,
	}, nil
}
