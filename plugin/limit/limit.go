package limit

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/miekg/dns"
)

type Inflight struct {
	cnt  int64
	max  int64
	next plugin.Handler
}

func (l *Inflight) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	defer func() {
		current := atomic.AddInt64(&l.cnt, -1)
		inflight.WithLabelValues(metrics.WithServer(ctx)).Set(float64(current))
	}()

	current := atomic.AddInt64(&l.cnt, 1)

	inflight.WithLabelValues(metrics.WithServer(ctx)).Set(float64(current))

	if current > l.max {
		return dns.RcodeServerFailure, plugin.Error(l.Name(), errors.New("max inflight request reached"))
	}

	return plugin.NextOrFailure(l.Name(), l.next, ctx, w, r)
}

func (l *Inflight) Name() string { return "inflight" }
