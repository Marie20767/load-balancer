package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

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
		counter: 1,
	}, nil
}

func (lb *LoadBalancer) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		if lb.counter == len(lb.config.Urls) {
			lb.counter = 1
		}

		urls := strings.Split(lb.config.Urls, ",")

		targetUrl, err := url.Parse(urls[lb.counter])
		if err != nil {
			return err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		req.Header.Set("X-Forwarded-For", c.RealIP())
		proxy.ServeHTTP(c.Response(), req)

		lb.counter++

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
