package provider

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func generateInputs() (Disk, ZDBInput, VMInput, QSFSInput, DeploymentArgs) {
	diskInput := Disk{
		Name:        "disk",
		Size:        0,
		Description: "disk-description",
	}

	zdbInput := ZDBInput{
		Name:        "zdb",
		Size:        1,
		Password:    "password",
		Public:      false,
		Description: "zdb-description",
		Mode:        "",
	}

	zlogInput := Zlog{
		Zmachine: "vm",
		Output:   "output",
	}

	vmInput := VMInput{
		Name:        "vm",
		NodeID:      1,
		Flist:       "flist",
		NetworkName: "network",
		CPU:         1,
		Memory:      1,
		PublicIP:    false,
		PublicIP6:   false,
		Planetary:   false,
		Description: "vm-description",
		GPUs:        []zos.GPU{"gpu"},
		RootfsSize:  0,
		Entrypoint:  "entrypoint",
		Mounts: []Mount{{
			Name:       diskInput.Name,
			MountPoint: fmt.Sprintf("/%s", diskInput.Name),
		}},
		Zlogs:   []Zlog{zlogInput},
		EnvVars: map[string]string{"key": "value"},
	}

	qsfsInput := QSFSInput{
		Name:                 "qsfs",
		Description:          "qsfs-description",
		Cache:                1,
		MinimalShards:        1,
		ExpectedShards:       1,
		RedundantGroups:      1,
		RedundantNodes:       1,
		MaxZDBDataDirSize:    1,
		EncryptionAlgorithm:  "algorithm",
		EncryptionKey:        "encryption",
		CompressionAlgorithm: "compression",
		Metadata: Metadata{
			EncryptionKey:       "encryption",
			Prefix:              "pre",
			EncryptionAlgorithm: "algorithm",
			Type:                "type",
			Backends: []Backend{{
				Address:   "address",
				Namespace: "namespace",
				Password:  "password",
			}},
		},
		Groups: []Group{{
			Backends: []Backend{{
				Address:   "address",
				Namespace: "namespace",
				Password:  "password",
			}},
		}},
	}

	deploymentInput := DeploymentArgs{
		NodeID:           1,
		Name:             "name",
		NetworkName:      "network",
		SolutionType:     "solution",
		SolutionProvider: 1,
		Disks:            []Disk{diskInput},
		ZdbsInputs:       []ZDBInput{zdbInput},
		VmsInputs:        []VMInput{vmInput},
		QSFSInputs:       []QSFSInput{qsfsInput},
	}

	return diskInput, zdbInput, vmInput, qsfsInput, deploymentInput
}

func TestDeploymentParser(t *testing.T) {
	diskInput, zdbInput, vmInput, qsfsInput, deploymentInput := generateInputs()

	t.Run("parsing input success", func(t *testing.T) {
		deployment, err := parseInputToDeployment(deploymentInput)
		assert.NoError(t, err)
		assert.Equal(t, deployment.NodeID, uint32(deploymentInput.NodeID.(int)))
		assert.Equal(t, deployment.Name, deploymentInput.Name)
		assert.Equal(t, deployment.NetworkName, deploymentInput.NetworkName)
		assert.Equal(t, deployment.SolutionType, deploymentInput.SolutionType)
		assert.Equal(t, *deployment.SolutionProvider, uint64(deploymentInput.SolutionProvider))
		assert.Equal(t, deployment.Disks[0].Name, diskInput.Name)
		assert.Equal(t, deployment.Zdbs[0].Name, zdbInput.Name)
		assert.Equal(t, deployment.Vms[0].Name, vmInput.Name)
		assert.Equal(t, deployment.QSFS[0].Name, qsfsInput.Name)
	})

	t.Run("parsing input failed: wrong node id type", func(t *testing.T) {
		deploymentInput.NodeID = ""
		_, err := parseInputToDeployment(deploymentInput)
		assert.Error(t, err)
		deploymentInput.NodeID = 1
	})

	t.Run("parsing and update deployment success", func(t *testing.T) {
		deployment, err := parseInputToDeployment(deploymentInput)
		assert.NoError(t, err)

		deployment.NodeDeploymentID = map[uint32]uint64{1: 1}
		state := parseDeploymentToState(deployment)
		assert.Equal(t, deployment.NodeDeploymentID[1], uint64(state.NodeDeploymentID["1"]))

		err = updateDeploymentFromState(&deployment, state)
		assert.NoError(t, err)
	})

	t.Run("parsing and update deployment failed: wrong node id type", func(t *testing.T) {
		deployment, err := parseInputToDeployment(deploymentInput)
		assert.NoError(t, err)

		state := parseDeploymentToState(deployment)
		state.NodeDeploymentID = map[string]int64{"": 1}
		err = updateDeploymentFromState(&deployment, state)
		assert.Error(t, err)
	})
}
