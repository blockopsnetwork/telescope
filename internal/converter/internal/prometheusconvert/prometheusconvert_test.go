package prometheusconvert_test

import (
	"testing"

	"github.com/blockopsnetwork/telescope/internal/converter/internal/prometheusconvert"
	"github.com/blockopsnetwork/telescope/internal/converter/internal/test_common"
	_ "github.com/blockopsnetwork/telescope/internal/static/metrics/instance"
)

func TestConvert(t *testing.T) {
	test_common.TestDirectory(t, "testdata", ".yaml", true, []string{}, prometheusconvert.Convert)
}
