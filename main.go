package main

import (
	"log"
	"os"

	"github.com/Marie20767/load-balancer/internal"
	"github.com/Marie20767/load-balancer/internal/load_balancer/weighted_robin"
	"github.com/Marie20767/load-balancer/internal/utils"
	"github.com/labstack/echo/v4"
)

func run() error {
	e := echo.New()
	e.HTTPErrorHandler = utils.CustomErrHandler

	c, err := config.ParseEnv()

	if err != nil {
		return err
	}

	s, err := config.LoadConfig()

	if err != nil {
		return err
	}

	lb, err := weightedrobin.NewLoadBalancer(c.Port, s)

	if err != nil {
		return err
	}

	e.Any("/*", lb.Handle())

	return e.Start(":" + c.Port)
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}
