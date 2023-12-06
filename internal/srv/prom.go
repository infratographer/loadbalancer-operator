package srv

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const subsystem = "load_balancer_operator"

var numberLoadBalancersRequestedGauge = promauto.NewGauge(
	prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "load_balancers_count",
		Help:      "Total count of currently deployed load balancers",
	},
)
