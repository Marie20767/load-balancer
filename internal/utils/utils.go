package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

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

	if _, isUrlError := err.(*url.Error); isUrlError {
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

	s1_Url := os.Getenv("SERVER_1_Url")
	s2_Url := os.Getenv("SERVER_2_Url")
	s3_Url := os.Getenv("SERVER_3_Url")

	if s1_Url == "" || s2_Url == "" || s3_Url == "" {
		return nil, errors.New("server_url environment variables not set")
	}

	s1_port := os.Getenv("SERVER_1_PORT")
	s2_port := os.Getenv("SERVER_2_PORT")
	s3_port := os.Getenv("SERVER_3_PORT")

	if s1_port == "" || s2_port == "" || s3_port == "" {
		return nil, errors.New("server_port environment variables not set")
	}

	urls := []string{
		fmt.Sprintf("%s:%s", s1_Url, s1_port),
		fmt.Sprintf("%s:%s", s2_Url, s2_port),
		fmt.Sprintf("%s:%s", s3_Url, s3_port),
	}
	urlList := strings.Join(urls, ",")

	return &config.Config{
		Port: port,
		Urls: urlList,
	}, nil
}
