//go:build !windows

package main

import (
	"testing"

	"github.com/blockopsnetwork/telescope/internal/cmd/integration-tests/common"
)

func TestScrapePromMetrics(t *testing.T) {
	common.MimirMetricsTest(t, common.PromDefaultMetrics, common.PromDefaultHistogramMetric, "scrape_prom_metrics")
}
