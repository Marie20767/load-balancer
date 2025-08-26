package roundrobinlb_test

import (
	"fmt"
	"testing"

	config "github.com/Marie20767/load-balancer/cmd/server/config"
	"github.com/Marie20767/load-balancer/internal/load_balancer/round_robin"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobinLb(t *testing.T) {
	p := "8080"
	u := "https://server_url_"

	t.Run("request should be forwarded to server 1", func(t *testing.T) {

		s := []config.Server{
			{URL: fmt.Sprintf("%s1.com", u), Weight: 1},
			{URL: fmt.Sprintf("%s2.com", u), Weight: 3},
			{URL: fmt.Sprintf("%s3.com", u), Weight: 1},
		}

		lb, err := roundrobinlb.NewLoadBalancer(p, s)
		assert.NoError(t, err)

		URL, err := lb.WeightedRoundRobin()
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

		lb, err := roundrobinlb.NewLoadBalancer(p, s)
		assert.NoError(t, err)

		_, err = lb.WeightedRoundRobin()
		assert.NoError(t, err)
		URL, err := lb.WeightedRoundRobin()
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

		lb, err := roundrobinlb.NewLoadBalancer(p, s)
		assert.NoError(t, err)

		_, err = lb.WeightedRoundRobin()
		assert.NoError(t, err)
		_, err = lb.WeightedRoundRobin()
		assert.NoError(t, err)
		URL, err := lb.WeightedRoundRobin()
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

		_, err := roundrobinlb.NewLoadBalancer(p, s)
		assert.ErrorIs(t, err, roundrobinlb.ErrWeight)
	})
}
