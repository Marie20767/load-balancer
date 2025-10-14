package weightedrobin

import (
	"errors"
	"net/http/httputil"
	"net/url"
	"sync"

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
	mu      sync.Mutex
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
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.counter > lb.weights {
		lb.counter = 1
	}

	sum := 0

	for _, server := range lb.servers {
		sum += server.Weight

		if lb.counter <= sum {
			lb.counter++
			u, err := url.Parse(server.URL)
			if err != nil {
				return nil, err
			}
			return u, nil
		}
	}

	return nil, ErrNoServer
}
