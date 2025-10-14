package consistenthashing

import (
	"errors"
	"hash/crc32"
	"math"
	"net/http/httputil"
	"net/url"

	"github.com/Marie20767/load-balancer/internal/loadbalancer/consistenthashing/config"
	"github.com/labstack/echo/v4"
)

var (
	ErrNoServer = errors.New("no server selected")
)

type LoadBalancer struct {
	servers []config.Server
	port    string
}

func NewLoadBalancer(port string, servers []config.Server) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		port:    port,
	}
}

func (lb *LoadBalancer) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ip := c.RealIP() // for testing different ips passed in X-Forwarded-For header, not for prod

		targetURL, err := lb.PickServer(ip)

		if err != nil {
			return err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		req.Header.Set("X-Forwarded-For", c.RealIP())
		proxy.ServeHTTP(c.Response(), req)

		return nil
	}
}

func (lb *LoadBalancer) PickServer(ip string) (*url.URL, error) {
	hash := HashInRange(ip)

	for _, s := range lb.servers {
		if hash < s.Position {
			return url.Parse(s.URL)
		}
	}

	return nil, ErrNoServer
}

func HashInRange(key string) float32 {
	hash := crc32.ChecksumIEEE([]byte(key))

	return (float32(hash) / float32(math.MaxUint32)) * config.HashRingValue
}
