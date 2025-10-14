package weightedrobin

import (
	"errors"
	"net/http/httputil"
	"net/url"

	"github.com/Marie20767/load-balancer/internal/loadbalancer/weightedrobin/config"
	"github.com/labstack/echo/v4"
)

var (
	ErrWeight   = errors.New("server weight cannot be 0")
	ErrNoServer = errors.New("no server selected")
)

type LoadBalancer struct {
	servers []config.Server
	weights int
	port    string
	counter int
}

func NewLoadBalancer(port string, servers []config.Server) (*LoadBalancer, error) {
	totalWeights := 0

	for _, server := range servers {
		if server.Weight == 0 {
			return nil, ErrWeight
		}
		totalWeights += server.Weight
	}

	return &LoadBalancer{
		port:    port,
		servers: servers,
		weights: totalWeights,
		counter: 1,
	}, nil
}

func (lb *LoadBalancer) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		targetURL, err := lb.PickServer()

		if err != nil {
			return err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		req.Header.Set("X-Forwarded-For", c.RealIP())
		proxy.ServeHTTP(c.Response(), req)

		return nil
	}
}

func (lb *LoadBalancer) PickServer() (*url.URL, error) {
	if lb.counter > lb.weights {
		lb.counter = 1
	}

	prevWeight := 0

	for i, server := range lb.servers {
		if i > 0 {
			prevWeight += lb.servers[i-1].Weight
		}

		if lb.counter <= prevWeight+server.Weight {
			lb.counter++
			return url.Parse(server.URL)
		}
	}

	return nil, ErrNoServer
}
