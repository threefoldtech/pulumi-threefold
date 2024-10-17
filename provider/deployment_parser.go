package provider

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/zos"
)

// Disk represents the disk struct
type Disk struct {
	Name        string `pulumi:"name"`
	Size        int    `pulumi:"size"`
	Description string `pulumi:"description,optional"`
}

// Mount represents mounting of disks
type Mount struct {
	Name       string `pulumi:"name"`
	MountPoint string `pulumi:"mount_point"`
}

// Zlog logger struct
type Zlog struct {
	Zmachine string `pulumi:"zmachine"`
	Output   string `pulumi:"output"`
}

// ZDBInput is the ZDB workload inputs struct
type ZDBInput struct {
	Name        string `pulumi:"name"`
	Size        int    `pulumi:"size"`
	Password    string `pulumi:"password"`
	Public      bool   `pulumi:"public,optional"`
	Description string `pulumi:"description,optional"`
	Mode        string `pulumi:"mode,optional"`
}

// ZDBComputed is the ZDB workload Computed struct
type ZDBComputed struct {
	IPs       []string `pulumi:"ips"`
	Port      int32    `pulumi:"port"`
	Namespace string   `pulumi:"namespace"`
}

// VMInput is a virtual machine inputs struct
type VMInput struct {
	Name           string            `pulumi:"name"`
	NodeID         interface{}       `pulumi:"node_id"`
	Flist          string            `pulumi:"flist"`
	FlistChecksum  string            `pulumi:"flist_checksum,optional"`
	NetworkName    string            `pulumi:"network_name"`
	CPU            int               `pulumi:"cpu"`
	Memory         int               `pulumi:"memory"`
	PublicIP       bool              `pulumi:"public_ip,optional"`
	PublicIP6      bool              `pulumi:"public_ip6,optional"`
	Planetary      bool              `pulumi:"planetary,optional"`
	Mycelium       bool              `pulumi:"mycelium,optional"`
	MyceliumIPSeed string            `pulumi:"mycelium_ip_seed,optional"`
	Description    string            `pulumi:"description,optional"`
	GPUs           []zos.GPU         `pulumi:"gpus,optional"`
	RootfsSize     int               `pulumi:"rootfs_size,optional"`
	EntryPoint     string            `pulumi:"entrypoint,optional"`
	Mounts         []Mount           `pulumi:"mounts,optional"`
	Zlogs          []Zlog            `pulumi:"zlogs,optional"`
	EnvVars        map[string]string `pulumi:"env_vars,optional"`
}

// VMComputed is a virtual machine computed struct
type VMComputed struct {
	MyceliumIPSeed string `pulumi:"mycelium_ip_seed"`
	ComputedIP     string `pulumi:"computed_ip"`
	ComputedIP6    string `pulumi:"computed_ip6"`
	PlanetaryIP    string `pulumi:"planetary_ip"`
	MyceliumIP     string `pulumi:"mycelium_ip"`
	ConsoleURL     string `pulumi:"console_url"`
	IP             string `pulumi:"ip,optional"`
}

// QSFSInput is the QSFS input struct
type QSFSInput struct {
	Name                 string   `pulumi:"name"`
	Description          string   `pulumi:"description,optional"`
	Cache                int      `pulumi:"cache"`
	MinimalShards        int32    `pulumi:"minimal_shards"`
	ExpectedShards       int32    `pulumi:"expected_shards"`
	RedundantGroups      int32    `pulumi:"redundant_groups"`
	RedundantNodes       int32    `pulumi:"redundant_nodes"`
	MaxZDBDataDirSize    int32    `pulumi:"max_zdb_data_dir_size"`
	EncryptionAlgorithm  string   `pulumi:"encryption_algorithm,optional"`
	EncryptionKey        string   `pulumi:"encryption_key"`
	CompressionAlgorithm string   `pulumi:"compression_algorithm,optional"`
	Metadata             Metadata `pulumi:"metadata"`
	Groups               []Group  `pulumi:"groups"`
}

// QSFSComputed is the QSFS Computed struct
type QSFSComputed struct {
	MetricsEndpoint string `pulumi:"metrics_endpoint"`
}

// Metadata for QSFS
type Metadata struct {
	EncryptionKey       string    `pulumi:"encryption_key"`
	Prefix              string    `pulumi:"prefix"`
	EncryptionAlgorithm string    `pulumi:"encryption_algorithm,optional"`
	Type                string    `pulumi:"type,optional"`
	Backends            []Backend `pulumi:"backends,optional"`
}

// Backend is a zos backend
type Backend struct {
	Address   string `pulumi:"address"`
	Namespace string `pulumi:"namespace"`
	Password  string `pulumi:"password"`
}

// Group is a zos group
type Group struct {
	Backends []Backend `pulumi:"backends,optional"`
}

func parseBackendsToState(backends workloads.Backends) []Backend {
	var bs []Backend
	for _, backend := range backends {
		bs = append(bs, Backend{
			Address:   backend.Address,
			Namespace: backend.Namespace,
			Password:  backend.Password,
		})
	}
	return bs
}

func parseInputsToBackends(backends []Backend) workloads.Backends {
	var bs workloads.Backends
	for _, backend := range backends {
		bs = append(bs, workloads.Backend{
			Address:   backend.Address,
			Namespace: backend.Namespace,
			Password:  backend.Password,
		})
	}
	return bs
}

func parseInputToDeployment(deploymentArgs DeploymentArgs, light bool) (workloads.Deployment, error) {
	var solutionProvider *uint64
	if deploymentArgs.SolutionProvider != 0 {
		solutionProviderUint := uint64(deploymentArgs.SolutionProvider)
		solutionProvider = &solutionProviderUint
	}

	nodeID, err := strconv.Atoi(fmt.Sprint(deploymentArgs.NodeID))
	if err != nil {
		return workloads.Deployment{}, err
	}

	var vms []workloads.VM
	var vmsLight []workloads.VMLight
	for _, vm := range deploymentArgs.VmsInputs {
		var mounts []workloads.Mount
		for _, mount := range vm.Mounts {
			mounts = append(mounts, workloads.Mount{
				Name:       mount.Name,
				MountPoint: mount.MountPoint,
			})
		}

		var zlogs []workloads.Zlog
		for _, zlog := range vm.Zlogs {
			zlogs = append(zlogs, workloads.Zlog{
				Zmachine: zlog.Zmachine,
				Output:   zlog.Output,
			})
		}

		// for tests
		sshKey := os.Getenv("SSH_KEY")
		if sshKey != "" {
			vm.EnvVars["SSH_KEY"] = sshKey
		}

		myceliumIPSeed, err := hex.DecodeString(vm.MyceliumIPSeed)
		if err != nil {
			return workloads.Deployment{}, err
		}

		if vm.Mycelium && len(vm.MyceliumIPSeed) == 0 {
			myceliumIPSeed, err = workloads.RandomMyceliumIPSeed()
			if err != nil {
				return workloads.Deployment{}, err
			}
		}

		if light {
			vmsLight = append(vmsLight, workloads.VMLight{
				Name:           vm.Name,
				NodeID:         uint32(nodeID),
				Flist:          vm.Flist,
				FlistChecksum:  vm.FlistChecksum,
				NetworkName:    vm.NetworkName,
				MyceliumIPSeed: myceliumIPSeed,
				Description:    vm.Description,
				GPUs:           vm.GPUs,
				CPU:            uint8(vm.CPU),
				MemoryMB:       uint64(vm.Memory),
				RootfsSizeMB:   uint64(vm.RootfsSize),
				Entrypoint:     vm.EntryPoint,
				Mounts:         mounts,
				Zlogs:          zlogs,
				EnvVars:        vm.EnvVars,
			})
			continue
		}

		vms = append(vms, workloads.VM{
			Name:           vm.Name,
			NodeID:         uint32(nodeID),
			Flist:          vm.Flist,
			FlistChecksum:  vm.FlistChecksum,
			NetworkName:    vm.NetworkName,
			PublicIP:       vm.PublicIP,
			PublicIP6:      vm.PublicIP6,
			Planetary:      vm.Planetary,
			MyceliumIPSeed: myceliumIPSeed,
			Description:    vm.Description,
			GPUs:           vm.GPUs,
			CPU:            uint8(vm.CPU),
			MemoryMB:       uint64(vm.Memory),
			RootfsSizeMB:   uint64(vm.RootfsSize),
			Entrypoint:     vm.EntryPoint,
			Mounts:         mounts,
			Zlogs:          zlogs,
			EnvVars:        vm.EnvVars,
		})
	}

	var disks []workloads.Disk
	for _, disk := range deploymentArgs.Disks {
		disks = append(disks, workloads.Disk{
			Name:        disk.Name,
			SizeGB:      uint64(disk.Size),
			Description: disk.Description,
		})
	}

	var zdbs []workloads.ZDB
	for _, zdb := range deploymentArgs.ZdbsInputs {
		zdbs = append(zdbs, workloads.ZDB{
			Name:        zdb.Name,
			SizeGB:      uint64(zdb.Size),
			Password:    zdb.Password,
			Public:      zdb.Public,
			Description: zdb.Description,
			Mode:        zdb.Mode,
		})
	}

	var qsfss []workloads.QSFS
	for _, qsfs := range deploymentArgs.QSFSInputs {
		var groups []workloads.Group
		for _, group := range qsfs.Groups {
			groups = append(groups, workloads.Group{
				Backends: parseInputsToBackends(group.Backends),
			})
		}

		qsfss = append(qsfss, workloads.QSFS{
			Name:                 qsfs.Name,
			Description:          qsfs.Description,
			Cache:                qsfs.Cache,
			MinimalShards:        uint32(qsfs.MinimalShards),
			ExpectedShards:       uint32(qsfs.ExpectedShards),
			RedundantGroups:      uint32(qsfs.RedundantGroups),
			RedundantNodes:       uint32(qsfs.RedundantNodes),
			MaxZDBDataDirSize:    uint32(qsfs.MaxZDBDataDirSize),
			EncryptionAlgorithm:  qsfs.EncryptionAlgorithm,
			EncryptionKey:        qsfs.EncryptionKey,
			CompressionAlgorithm: qsfs.CompressionAlgorithm,
			Metadata: workloads.Metadata{
				EncryptionKey:       qsfs.Metadata.EncryptionKey,
				Prefix:              qsfs.Metadata.Prefix,
				EncryptionAlgorithm: qsfs.Metadata.EncryptionAlgorithm,
				Type:                qsfs.Metadata.Type,
				Backends:            parseInputsToBackends(qsfs.Metadata.Backends),
			},
			Groups: groups,
		})
	}

	return workloads.Deployment{
		NodeID:           uint32(nodeID),
		Name:             deploymentArgs.Name,
		SolutionType:     deploymentArgs.SolutionType,
		SolutionProvider: solutionProvider,
		NetworkName:      deploymentArgs.NetworkName,
		Disks:            disks,
		Zdbs:             zdbs,
		Vms:              vms,
		VmsLight:         vmsLight,
		QSFS:             qsfss,
	}, nil
}

func parseDeploymentToState(deployment workloads.Deployment) DeploymentState {
	var solutionProvider int64
	if deployment.SolutionProvider != nil {
		solutionProvider = int64(*deployment.SolutionProvider)
	}

	var disks []Disk
	for _, disk := range deployment.Disks {
		disks = append(disks, Disk{
			Name:        disk.Name,
			Size:        int(disk.SizeGB),
			Description: disk.Description,
		})
	}

	var zdbs []ZDBInput
	for _, zdb := range deployment.Zdbs {
		zdbs = append(zdbs, ZDBInput{
			Name:        zdb.Name,
			Size:        int(zdb.SizeGB),
			Password:    zdb.Password,
			Public:      zdb.Public,
			Description: zdb.Description,
			Mode:        zdb.Mode,
		})
	}

	var vms []VMInput
	for _, vm := range deployment.Vms {
		var zlogs []Zlog
		for _, zlog := range vm.Zlogs {
			zlogs = append(zlogs, Zlog{
				Zmachine: zlog.Zmachine,
				Output:   zlog.Output,
			})
		}

		var mounts []Mount
		for _, mount := range vm.Mounts {
			mounts = append(mounts, Mount{
				Name:       mount.Name,
				MountPoint: mount.MountPoint,
			})
		}

		if vm.NodeID == 0 {
			vm.NodeID = deployment.NodeID
		}

		vms = append(vms, VMInput{
			Name:           vm.Name,
			NodeID:         vm.NodeID,
			Flist:          vm.Flist,
			FlistChecksum:  vm.FlistChecksum,
			NetworkName:    vm.NetworkName,
			PublicIP:       vm.PublicIP,
			PublicIP6:      vm.PublicIP6,
			Planetary:      vm.Planetary,
			MyceliumIPSeed: hex.EncodeToString(vm.MyceliumIPSeed),
			Description:    vm.Description,
			GPUs:           vm.GPUs,
			CPU:            int(vm.CPU),
			Memory:         int(vm.MemoryMB),
			RootfsSize:     int(vm.RootfsSizeMB),
			EntryPoint:     vm.Entrypoint,
			Mounts:         mounts,
			Zlogs:          zlogs,
			EnvVars:        vm.EnvVars,
		})
	}

	for _, vm := range deployment.VmsLight {
		var zlogs []Zlog
		for _, zlog := range vm.Zlogs {
			zlogs = append(zlogs, Zlog{
				Zmachine: zlog.Zmachine,
				Output:   zlog.Output,
			})
		}

		var mounts []Mount
		for _, mount := range vm.Mounts {
			mounts = append(mounts, Mount{
				Name:       mount.Name,
				MountPoint: mount.MountPoint,
			})
		}

		if vm.NodeID == 0 {
			vm.NodeID = deployment.NodeID
		}

		vms = append(vms, VMInput{
			Name:           vm.Name,
			NodeID:         vm.NodeID,
			Flist:          vm.Flist,
			FlistChecksum:  vm.FlistChecksum,
			NetworkName:    vm.NetworkName,
			MyceliumIPSeed: hex.EncodeToString(vm.MyceliumIPSeed),
			Description:    vm.Description,
			GPUs:           vm.GPUs,
			CPU:            int(vm.CPU),
			Memory:         int(vm.MemoryMB),
			RootfsSize:     int(vm.RootfsSizeMB),
			EntryPoint:     vm.Entrypoint,
			Mounts:         mounts,
			Zlogs:          zlogs,
			EnvVars:        vm.EnvVars,
		})
	}

	var qsfss []QSFSInput
	for _, qsfs := range deployment.QSFS {
		var groups []Group
		for _, group := range qsfs.Groups {
			groups = append(groups, Group{
				Backends: parseBackendsToState(group.Backends),
			})
		}

		qsfss = append(qsfss, QSFSInput{
			Name:                 qsfs.Name,
			Description:          qsfs.Description,
			Cache:                qsfs.Cache,
			MinimalShards:        int32(qsfs.MinimalShards),
			ExpectedShards:       int32(qsfs.ExpectedShards),
			RedundantGroups:      int32(qsfs.RedundantGroups),
			RedundantNodes:       int32(qsfs.RedundantNodes),
			MaxZDBDataDirSize:    int32(qsfs.MaxZDBDataDirSize),
			EncryptionAlgorithm:  qsfs.EncryptionAlgorithm,
			EncryptionKey:        qsfs.EncryptionKey,
			CompressionAlgorithm: qsfs.CompressionAlgorithm,
			Metadata: Metadata{
				EncryptionKey:       qsfs.Metadata.EncryptionKey,
				Prefix:              qsfs.Metadata.Prefix,
				EncryptionAlgorithm: qsfs.Metadata.EncryptionAlgorithm,
				Type:                qsfs.Metadata.Type,
				Backends:            parseBackendsToState(qsfs.Metadata.Backends),
			},
			Groups: groups,
		})
	}

	stateArgs := DeploymentArgs{
		NodeID:           int(deployment.NodeID),
		Name:             deployment.Name,
		SolutionType:     deployment.SolutionType,
		SolutionProvider: solutionProvider,
		NetworkName:      deployment.NetworkName,
		Disks:            disks,
		ZdbsInputs:       zdbs,
		VmsInputs:        vms,
		QSFSInputs:       qsfss,
	}

	nodeDeploymentID := make(map[string]int64)
	for nodeID, deploymentID := range deployment.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(nodeID)] = int64(deploymentID)
	}

	zdbsComputed := make([]ZDBComputed, 0)
	for _, zdb := range deployment.Zdbs {
		zdbsComputed = append(zdbsComputed, ZDBComputed{
			IPs:       zdb.IPs,
			Port:      int32(zdb.Port),
			Namespace: zdb.Namespace,
		})
	}

	vmsComputed := make([]VMComputed, 0)
	for _, vm := range deployment.Vms {
		vmsComputed = append(vmsComputed, VMComputed{
			MyceliumIPSeed: hex.EncodeToString(vm.MyceliumIPSeed),
			ComputedIP:     vm.ComputedIP,
			ComputedIP6:    vm.ComputedIP6,
			MyceliumIP:     vm.MyceliumIP,
			PlanetaryIP:    vm.PlanetaryIP,
			ConsoleURL:     vm.ConsoleURL,
			IP:             vm.IP,
		})
	}

	for _, vm := range deployment.VmsLight {
		vmsComputed = append(vmsComputed, VMComputed{
			MyceliumIPSeed: hex.EncodeToString(vm.MyceliumIPSeed),
			MyceliumIP:     vm.MyceliumIP,
			ConsoleURL:     vm.ConsoleURL,
			IP:             vm.IP,
		})
	}

	qsfsComputed := make([]QSFSComputed, 0)
	for _, qsfs := range deployment.QSFS {
		qsfsComputed = append(qsfsComputed, QSFSComputed{
			MetricsEndpoint: qsfs.MetricsEndpoint,
		})
	}

	state := DeploymentState{
		DeploymentArgs:   stateArgs,
		NodeDeploymentID: nodeDeploymentID,
		ContractID:       int64(deployment.ContractID),
		IPrange:          deployment.IPrange,
		ZdbsComputed:     zdbsComputed,
		VmsComputed:      vmsComputed,
		QsfsComputed:     qsfsComputed,
	}

	return state
}

func updateDeploymentFromState(deployment *workloads.Deployment, state DeploymentState, light bool) error {
	nodeDeploymentID := make(map[uint32]uint64)
	for nodeID, deploymentID := range state.NodeDeploymentID {
		node, err := strconv.ParseUint(nodeID, 10, 32)
		if err != nil {
			return err
		}
		nodeDeploymentID[uint32(node)] = uint64(deploymentID)
	}

	for i, zdb := range state.ZdbsComputed {
		if len(deployment.Zdbs) > i {
			deployment.Zdbs[i].IPs = zdb.IPs
			deployment.Zdbs[i].Port = uint32(zdb.Port)
			deployment.Zdbs[i].Namespace = zdb.Namespace
		}
	}

	for i, vm := range state.VmsComputed {
		if len(deployment.Vms) > i || len(deployment.VmsLight) > i {
			myceliumIPSeed, err := hex.DecodeString(vm.MyceliumIPSeed)
			if err != nil {
				return err
			}

			if light {
				deployment.VmsLight[i].MyceliumIPSeed = myceliumIPSeed
				deployment.VmsLight[i].MyceliumIP = vm.MyceliumIP
				deployment.VmsLight[i].ConsoleURL = vm.ConsoleURL
				deployment.VmsLight[i].IP = vm.IP
				continue
			}

			deployment.Vms[i].MyceliumIPSeed = myceliumIPSeed
			deployment.Vms[i].ComputedIP = vm.ComputedIP
			deployment.Vms[i].ComputedIP6 = vm.ComputedIP6
			deployment.Vms[i].PlanetaryIP = vm.PlanetaryIP
			deployment.Vms[i].MyceliumIP = vm.MyceliumIP
			deployment.Vms[i].ConsoleURL = vm.ConsoleURL
			deployment.Vms[i].IP = vm.IP
		}
	}

	for i, qsfs := range state.QsfsComputed {
		if len(deployment.QSFS) > i {
			deployment.QSFS[i].MetricsEndpoint = qsfs.MetricsEndpoint
		}
	}

	deployment.NodeDeploymentID = nodeDeploymentID
	deployment.ContractID = uint64(state.ContractID)
	deployment.IPrange = state.IPrange

	return nil
}
