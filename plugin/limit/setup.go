package limit

import (
	"errors"
	"strconv"
	"sync"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/caddyserver/caddy"
)

var (
	log  = clog.NewWithPlugin("limit")
	once sync.Once
)

func init() { plugin.Register("limit", setup) }

func setup(c *caddy.Controller) error {
	max, err := parse(c)
	if err != nil {
		return plugin.Error("limit", err)
	}

	c.OnStartup(func() error {
		once.Do(func() { metrics.MustRegister(c, inflight) })
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return &Inflight{max: max, next: next}
	})

	return nil
}

func parse(c *caddy.Controller) (int64, error) {
	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			return 0, errors.New("unknown limits")
		case 1:
			return 0, errors.New("unknown limits")
		case 2:
			if args[0] == "requests" {
				return strconv.ParseInt(args[1], 10, 64)
			}
		}
	}

	return 0, c.ArgErr()
}
