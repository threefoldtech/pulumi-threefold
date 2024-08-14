package test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	mnemonic := os.Getenv("MNEMONIC")
	assert.NotEmpty(t, mnemonic)

	network := os.Getenv("NETWORK")
	if network == "" {
		network = devNetwork
	}

	examplesDir := os.Getenv("EXAMPLES")
	if examplesDir == "" {
		examplesDir = examplesTestDir
	}

	cwd, _ := os.Getwd()

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Quick:            true,
		SkipRefresh:      true,
		DestroyOnCleanup: true,
		Dir:              path.Join(cwd, fmt.Sprintf("%s/network", examplesDir)),
		Config: map[string]string{
			"MNEMONIC": mnemonic,
			"NETWORK":  network,
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			for _, res := range stack.Deployment.Resources {
				if res.Type == "threefold:Network" {
					assert.NotEmpty(t, res.Outputs["node_deployment_id"])
					assert.NotEmpty(t, res.Outputs["nodes_ip_range"])
				}
			}
		},
	})
}
