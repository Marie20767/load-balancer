package main

import (
	"log"
	"os"

	"github.com/Marie20767/load-balancer/internal/loadbalancer/consistenthashing"
	"github.com/Marie20767/load-balancer/internal/loadbalancer/consistenthashing/config"

	"github.com/Marie20767/load-balancer/internal/utils"
	"github.com/labstack/echo/v4"
)

func run() error {
	e := echo.New()
	e.HTTPErrorHandler = utils.CustomErrHandler

	c, err := utils.ParseEnv()

	if err != nil {
		return err
	}

	s, err := config.LoadConfig()

	if err != nil {
		return err
	}

	lb := consistenthashing.NewLoadBalancer(c.Port, s)

	e.Any("/*", lb.Handle())

	return e.Start(":" + c.Port)
}

func main() {
	if err := run(); err != nil {
		log.Println("server closed: ", err)
		os.Exit(1)
	}
}
