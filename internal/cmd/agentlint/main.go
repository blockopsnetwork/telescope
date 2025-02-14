// Command agentlint provides custom linting utilities for the grafana/agent
// repo.
package main

import (
	"github.com/blockopsnetwork/telescope/internal/cmd/agentlint/internal/findcomponents"
	"github.com/blockopsnetwork/telescope/internal/cmd/agentlint/internal/rivertags"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		findcomponents.Analyzer,
		rivertags.Analyzer,
	)
}
