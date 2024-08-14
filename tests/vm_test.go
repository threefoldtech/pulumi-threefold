package test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
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

	publicKey, privateKey, err := generateSSHKeyPair()
	assert.NoError(t, err)

	err = os.Setenv("SSH_KEY", publicKey)
	assert.NoError(t, err)

	cwd, _ := os.Getwd()

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Quick:            true,
		SkipRefresh:      true,
		DestroyOnCleanup: true,
		Dir:              path.Join(cwd, fmt.Sprintf("%s/virtual_machine", examplesDir)),
		Config: map[string]string{
			"MNEMONIC": mnemonic,
			"NETWORK":  network,
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			for _, res := range stack.Deployment.Resources {
				if res.Type == "threefold:Deployment" {
					assert.NotEmpty(t, res.Outputs["node_deployment_id"])

					vmsComputed := res.Outputs["vms_computed"].([]interface{})[0].(map[string]interface{})
					vms := res.Outputs["vms"].([]interface{})[0].(map[string]interface{})
					disks := res.Outputs["disks"].([]interface{})[0].(map[string]interface{})
					mounts := vms["mounts"].([]interface{})[0].(map[string]interface{})

					yggIP := vmsComputed["planetary_ip"].(string)
					mountPoint := mounts["mount_point"].(string)
					diskSize := disks["size"].(float64)

					// testing connection
					ok := testConnection(yggIP, "22")
					assert.True(t, ok)

					// Check that disk has been mounted successfully
					output, err := remoteRun("root", yggIP, fmt.Sprintf("df -h | grep -w %s", mountPoint), privateKey)
					assert.NoError(t, err)
					assert.Contains(t, output, fmt.Sprintf("%v.0G", diskSize))
				}
			}
		},
	})
}
