package main

import (
	"errors"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/Marie20767/load-balancer/internal/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Spec struct {
	Url     string
	S1_port string
	S2_port string
	S3_port string
}

func loadBalancer(spec *Spec) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		target := fmt.Sprintf("%s:%s", spec.Url, spec.S1_port)
		targetURL, err := url.Parse(target)
		if err != nil {
			return err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		req.Header.Set("X-Forwarded-For", c.RealIP())
		proxy.ServeHTTP(c.Response(), req)

		return nil
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	port := os.Getenv("PORT")
	URL := os.Getenv("SERVER_URL")
	s1_port := os.Getenv("SERVER_1_PORT")
	s2_port := os.Getenv("SERVER_2_PORT")
	s3_port := os.Getenv("SERVER_3_PORT")

	spec := &Spec{
		Url:     URL,
		S1_port: s1_port,
		S2_port: s2_port,
		S3_port: s3_port,
	}

	if spec.Url == "" || spec.S1_port == "" || spec.S2_port == "" || spec.S3_port == "" || port == "" {
		return errors.New("not all environment variables are set")
	}

	e := echo.New()
	e.HTTPErrorHandler = utils.CustomErrHandler
	e.Any("/*", loadBalancer(spec))

	return e.Start(":" + port)
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}
