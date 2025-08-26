package main

import (
	"log"
	"os"

	serverConfig "github.com/Marie20767/load-balancer/cmd/server/config"
	config "github.com/Marie20767/load-balancer/internal"
	"github.com/Marie20767/load-balancer/internal/load_balancer/round_robin"
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

	s, err := serverConfig.LoadConfig()

	if err != nil {
		return err
	}

	lb, err := roundrobinlb.NewLoadBalancer(c.Port, s)

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
