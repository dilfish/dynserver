package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var badDomainNameCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "bad_domain_name_count",
		Help: "count of domain name does not match",
	},
)

var requestDurationMs = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "request_duration_ms",
		Help: "milliseconds a request cost",
	},
)

var fileSizeServedBytes = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "file_size_served_bytes",
		Help: "served file in bytes",
	},
)