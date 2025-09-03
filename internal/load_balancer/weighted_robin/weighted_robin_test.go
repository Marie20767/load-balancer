package weightedrobin_test

import (
	"fmt"
	"testing"

	"github.com/Marie20767/load-balancer/internal/load_balancer/weighted_robin"
	"github.com/Marie20767/load-balancer/internal/load_balancer/weighted_robin/config"

	"github.com/stretchr/testify/assert"
)

func TestWeightedRobin(t *testing.T) {
	p := "8080"
	u := "https://server_url_"

	t.Run("request should be forwarded to server 1", func(t *testing.T) {

		s := []config.Server{
			{URL: fmt.Sprintf("%s1.com", u), Weight: 1},
			{URL: fmt.Sprintf("%s2.com", u), Weight: 3},
			{URL: fmt.Sprintf("%s3.com", u), Weight: 1},
		}

		lb, err := weightedrobin.NewLoadBalancer(p, s)
		assert.NoError(t, err)

		URL, err := lb.PickServer()
		expected := s[0].URL

		assert.NoError(t, err)
		assert.Equal(t, expected, URL.String())
	})

	t.Run("request should be forwarded to server 2", func(t *testing.T) {

		s := []config.Server{
			{URL: fmt.Sprintf("%s1.com", u), Weight: 1},
			{URL: fmt.Sprintf("%s2.com", u), Weight: 2},
			{URL: fmt.Sprintf("%s3.com", u), Weight: 1},
		}

		lb, err := weightedrobin.NewLoadBalancer(p, s)
		assert.NoError(t, err)

		_, err = lb.PickServer()
		assert.NoError(t, err)
		URL, err := lb.PickServer()
		expected := s[1].URL

		assert.NoError(t, err)
		assert.Equal(t, expected, URL.String())
	})

	t.Run("request should be forwarded to server 3", func(t *testing.T) {
		s := []config.Server{
			{URL: fmt.Sprintf("%s1.com", u), Weight: 1},
			{URL: fmt.Sprintf("%s2.com", u), Weight: 1},
			{URL: fmt.Sprintf("%s3.com", u), Weight: 1},
		}

		lb, err := weightedrobin.NewLoadBalancer(p, s)
		assert.NoError(t, err)

		_, err = lb.PickServer()
		assert.NoError(t, err)
		_, err = lb.PickServer()
		assert.NoError(t, err)
		URL, err := lb.PickServer()
		expected := s[2].URL

		assert.NoError(t, err)
		assert.Equal(t, expected, URL.String())
	})

	t.Run("cannot assign empty weight to server", func(t *testing.T) {
		s := []config.Server{
			{URL: fmt.Sprintf("%s1.com", u), Weight: 1},
			{URL: fmt.Sprintf("%s2.com", u), Weight: 0},
			{URL: fmt.Sprintf("%s3.com", u), Weight: 1},
		}

		_, err := weightedrobin.NewLoadBalancer(p, s)
		assert.ErrorIs(t, err, weightedrobin.ErrWeight)
	})
}
