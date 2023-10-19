package test

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestKubernetes(t *testing.T) {
	mnemonic := os.Getenv("MNEMONIC")
	assert.NotEmpty(t, mnemonic)

	network := os.Getenv("NETWORK")
	if network == "" {
		network = devNetwork
	}

	publicKey, privateKey, err := generateSSHKeyPair()
	assert.NoError(t, err)

	err = os.Setenv("SSH_KEY", publicKey)
	assert.NoError(t, err)

	cwd, _ := os.Getwd()

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Quick:            true,
		SkipRefresh:      true,
		DestroyOnCleanup: true,
		Dir:              path.Join(cwd, "..", "examples/kubernetes"),
		Config: map[string]string{
			"MNEMONIC": mnemonic,
			"NETWORK":  network,
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			for _, res := range stack.Deployment.Resources {
				if res.Type == "threefold:provider:Kubernetes" {
					assert.NotEmpty(t, res.Outputs["node_deployment_id"])

					yggIP := res.Outputs["master_computed"].(map[string]interface{})["ygg_ip"].(string)
					assert.NotEmpty(t, yggIP)
					AssertNodesAreReady(t, yggIP, privateKey, 3)
				}
			}
		},
	})
}

func AssertNodesAreReady(t *testing.T, masterYggIP, privateKey string, nodesNumber int) {
	t.Helper()

	time.Sleep(60 * time.Second)
	output, err := remoteRun(
		"root",
		masterYggIP,
		"export KUBECONFIG=/etc/rancher/k3s/k3s.yaml && kubectl get node",
		privateKey,
	)
	output = strings.TrimSpace(output)
	assert.Empty(t, err)

	numberOfReadyNodes := strings.Count(output, "Ready")
	assert.True(
		t,
		numberOfReadyNodes == nodesNumber,
		"number of ready nodes is not equal to number of nodes only %d nodes are ready",
		numberOfReadyNodes,
	)
}
