package promtailconvert_test

import (
	"testing"

	"github.com/blockopsnetwork/telescope/internal/converter/internal/promtailconvert"
	"github.com/blockopsnetwork/telescope/internal/converter/internal/test_common"
	_ "github.com/blockopsnetwork/telescope/internal/static/metrics/instance" // Imported to override default values via the init function.
)

func TestConvert(t *testing.T) {
	test_common.TestDirectory(t, "testdata", ".yaml", true, []string{}, promtailconvert.Convert)
}
