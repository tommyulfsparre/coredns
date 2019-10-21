package limit

import (
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
)

var inflight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: plugin.Namespace,
	Subsystem: "limit",
	Name:      "request_inflight",
	Help:      "Current number of inflight DNS requests.",
}, []string{"server"})
