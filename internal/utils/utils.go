package utils

import (
	"errors"
	"net/http"
	"net/url"
	"os"

	config "github.com/Marie20767/load-balancer/internal"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func CustomErrHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var code int
	var msg string

	if _, isURLError := err.(*url.Error); isURLError {
		code = http.StatusBadGateway
		msg = "Bad Gateway: " + err.Error()
	} else {
		code = http.StatusInternalServerError
		msg = err.Error()
	}

	c.String(code, msg)
}

func ParseEnv() (c *config.Config, err error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, errors.New("port environment variable not set")
	}

	s1_URL := os.Getenv("SERVER_1_URL")
	s2_URL := os.Getenv("SERVER_2_URL")
	s3_URL := os.Getenv("SERVER_3_URL")

	if s1_URL == "" || s2_URL == "" || s3_URL == "" {
		return nil, errors.New("server_url environment variables not set")
	}

	s1_port := os.Getenv("SERVER_1_PORT")
	s2_port := os.Getenv("SERVER_2_PORT")
	s3_port := os.Getenv("SERVER_3_PORT")

	if s1_port == "" || s2_port == "" || s3_port == "" {
		return nil, errors.New("server_port environment variables not set")
	}

	return &config.Config{
		Port:    port,
		S1_url:  s1_URL,
		S1_port: s1_port,
		S2_url:  s2_URL,
		S2_port: s2_port,
		S3_url:  s1_URL,
		S3_port: s1_URL,
	}, nil
}
