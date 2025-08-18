package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	config "github.com/Marie20767/load-balancer/internal"
	"github.com/Marie20767/load-balancer/internal/utils"
	"github.com/labstack/echo/v4"
)

type LoadBalancer struct {
	config  *config.Config
	counter int
}

func NewLoadBalancer() (lb *LoadBalancer, err error) {
	newConfig, err := utils.ParseEnv()

	if err != nil {
		return nil, err
	}

	return &LoadBalancer{
		config:  newConfig,
		counter: 0,
	}, nil
}

func (lb *LoadBalancer) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		target := fmt.Sprintf("%s:%s", lb.config.S1_url, lb.config.S1_port)
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
	e := echo.New()
	e.HTTPErrorHandler = utils.CustomErrHandler
	loadBalancer, err := NewLoadBalancer()

	if err != nil {
		return err
	}

	e.Any("/*", loadBalancer.Handle())

	return e.Start(":" + loadBalancer.config.Port)
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}
