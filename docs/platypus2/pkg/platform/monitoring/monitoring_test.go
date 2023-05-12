package monitoring

import (
	"os"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/lingoneks/pkg/platform/monitoring/metricsserver"
	"github.com/volvo-cars/lingoneks/pkg/platform/monitoring/promcrd"
	"github.com/volvo-cars/lingoneks/pkg/platform/monitoring/promstack"
)

func TestMonitoring(t *testing.T) {
	folders := []string{
		"out/1_promcrd",
		"out/2_metrics-server",
		"out/3_promstack",
	}
	for _, f := range folders {
		_ = os.RemoveAll(f)
	}

	pcrd := promcrd.New()
	if err := pcrd.Export(folders[0]); err != nil {
		tu.AssertNoError(t, err, "prometheus crd")
	}
	ms := metricsserver.New()
	if err := ms.Export(folders[1]); err != nil {
		tu.AssertNoError(t, err, "metrics-server")
	}

	ps := promstack.New()
	if err := ps.Export(folders[2]); err != nil {
		tu.AssertNoError(t, err, "prometheus stack")
	}
}
