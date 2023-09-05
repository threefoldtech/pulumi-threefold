package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

// Disk respresents the disk struct
type Disk struct {
	Name        string `pulumi:"name"`
	Size        int    `pulumi:"size"`
	Description string `pulumi:"description,optional"`
}

// Mount respresents mounting of disks
type Mount struct {
	DiskName   string `pulumi:"disk_name"`
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
	public      bool   `pulumi:"public,optional"`
	Description string `pulumi:"description,optional"`
	Mode        string `pulumi:"mode,optional"`
}

// ZDBComputed is the ZDB workload Computed struct
type ZDBComputed struct {
	IPs       []string `pulumi:"ips"`
	Port      uint32   `pulumi:"port"`
	Namespace string   `pulumi:"namespace"`
}

// VMInput is a virtual machine inputs struct
type VMInput struct {
	Name          string            `pulumi:"name"`
	Flist         string            `pulumi:"flist"`
	NetworkName   string            `pulumi:"network_name"`
	FlistChecksum string            `pulumi:"flist_checksum,optional"`
	PublicIP      bool              `pulumi:"publicip,optional"`
	PublicIP6     bool              `pulumi:"publicip6,optional"`
	Planetary     bool              `pulumi:"planetary,optional"`
	Corex         bool              `pulumi:"corex,optional"`
	Description   string            `pulumi:"description,optional"`
	GPUs          []zos.GPU         `pulumi:"gpus,optional"`
	CPU           int               `pulumi:"cpu,optional"`
	Memory        int               `pulumi:"memory,optional"`
	RootfsSize    int               `pulumi:"rootfs_size,optional"`
	Entrypoint    string            `pulumi:"entrypoint,optional"`
	Mounts        []Mount           `pulumi:"mounts,optional"`
	Zlogs         []Zlog            `pulumi:"zlogs,optional"`
	EnvVars       map[string]string `pulumi:"env_vars,optional"`
}

// VMComputed is a virtual machine computed struct
type VMComputed struct {
	ComputedIP  string `pulumi:"computedip"`
	ComputedIP6 string `pulumi:"computedip6"`
	YggIP       string `pulumi:"ygg_ip"`
	ConsoleURL  string `pulumi:"console_url"`
	IP          string `pulumi:"ip,optional"`
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

// Deployment controlling struct
type Deployment struct{}

// DeploymentArgs is defining what arguments it accepts
type DeploymentArgs struct {
	NodeID           int32       `pulumi:"node_id"`
	Name             string      `pulumi:"name,optional"`
	SolutionType     string      `pulumi:"solution_type,optional"`
	SolutionProvider *uint64     `pulumi:"solution_provider,optional"`
	NetworkName      string      `pulumi:"network_name,optional"`
	Disks            []Disk      `pulumi:"disks,optional"`
	ZdbsInputs       []ZDBInput  `pulumi:"zdbs_inputs,optional"`
	VmsInputs        []VMInput   `pulumi:"vms_inputs,optional"`
	QSFSInputs       []QSFSInput `pulumi:"qsfs_inputs,optional"`
}

// DeploymentState is describing the fields that exist on the created resource
type DeploymentState struct {
	DeploymentArgs

	NodeDeploymentID map[uint32]uint64
	ContractID       uint64
	IPrange          string
	ZdbsComputed     []ZDBComputed
	VmsComputed      []VMComputed
	QsfsComputed     []QSFSComputed
}

func convertDisks(disks []workloads.Disk) []Disk {
	result := make([]Disk, len(disks))
	for i, disk := range disks {
		result[i] = Disk{
			Name:        disk.Name,
			Size:        disk.SizeGB,
			Description: disk.Description,
		}
	}
	return result
}

func convertToWorkloadDisks(disks []Disk) []workloads.Disk {
	result := make([]workloads.Disk, len(disks))
	for i, disk := range disks {
		result[i] = workloads.Disk{
			Name:        disk.Name,
			SizeGB:      disk.Size,
			Description: disk.Description,
		}
	}
	return result
}

func convertZdbs(zdbs []workloads.ZDB) []ZDBInput {

	result := make([]ZDBInput, len(zdbs))

	for i, zdb := range zdbs {

		result[i] = ZDBInput{
			Name:        zdb.Name,
			Size:        zdb.Size,
			Password:    zdb.Password,
			public:      zdb.Public,
			Description: zdb.Description,
			Mode:        zdb.Mode,
		}
	}

	return result
}

func convertToWorkloadZdbs(zdbs []ZDBInput) []workloads.ZDB {

	result := make([]workloads.ZDB, len(zdbs))
	for i, zdb := range zdbs {

		result[i] = workloads.ZDB{
			Name:        zdb.Name,
			Size:        zdb.Size,
			Password:    zdb.Password,
			Public:      zdb.public,
			Description: zdb.Description,
			Mode:        zdb.Mode,
		}
	}
	return result
}

func convertZdbsComputed(zdbs []workloads.ZDB) []ZDBComputed {

	result := make([]ZDBComputed, len(zdbs))

	for i, zdb := range zdbs {

		result[i] = ZDBComputed{
			IPs:       zdb.IPs,
			Port:      zdb.Port,
			Namespace: zdb.Namespace,
		}
	}

	return result
}

func convertMounts(mounts []workloads.Mount) []Mount {
	result := make([]Mount, len(mounts))
	for i, mount := range mounts {
		result[i] = Mount{
			DiskName:   mount.DiskName,
			MountPoint: mount.MountPoint,
		}
	}
	return result
}

func convertToWorkloadMounts(mounts []Mount) []workloads.Mount {
	result := make([]workloads.Mount, len(mounts))
	for i, mount := range mounts {
		result[i] = workloads.Mount{
			DiskName:   mount.DiskName,
			MountPoint: mount.MountPoint,
		}
	}
	return result
}

func convertZlogs(zlogs []workloads.Zlog) []Zlog {
	result := make([]Zlog, len(zlogs))
	for i, zlog := range zlogs {
		result[i] = Zlog{
			Zmachine: zlog.Zmachine,
			Output:   zlog.Output,
		}
	}
	return result
}

func convertToWorkloadZlogs(zlogs []Zlog) []workloads.Zlog {
	result := make([]workloads.Zlog, len(zlogs))
	for i, zlog := range zlogs {
		result[i] = workloads.Zlog{
			Zmachine: zlog.Zmachine,
			Output:   zlog.Output,
		}
	}
	return result
}

func convertVMs(VMs []workloads.VM) []VMInput {

	result := make([]VMInput, len(VMs))

	for i, vm := range VMs {

		result[i] = VMInput{
			Name:          vm.Name,
			Flist:         vm.Flist,
			NetworkName:   vm.NetworkName,
			FlistChecksum: vm.FlistChecksum,
			PublicIP:      vm.PublicIP,
			PublicIP6:     vm.PublicIP6,
			Planetary:     vm.Planetary,
			Corex:         vm.Corex,
			Description:   vm.Description,
			GPUs:          vm.GPUs,
			CPU:           vm.CPU,
			Memory:        vm.Memory,
			RootfsSize:    vm.RootfsSize,
			Entrypoint:    vm.Entrypoint,
			Mounts:        convertMounts(vm.Mounts),
			Zlogs:         convertZlogs(vm.Zlogs),
			EnvVars:       vm.EnvVars,
		}
	}
	return result
}

func convertToWorkloadVMs(VMs []VMInput) []workloads.VM {

	result := make([]workloads.VM, len(VMs))

	for i, vm := range VMs {

		result[i] = workloads.VM{
			Name:          vm.Name,
			Flist:         vm.Flist,
			NetworkName:   vm.NetworkName,
			FlistChecksum: vm.FlistChecksum,
			PublicIP:      vm.PublicIP,
			PublicIP6:     vm.PublicIP6,
			Planetary:     vm.Planetary,
			Corex:         vm.Corex,
			Description:   vm.Description,
			GPUs:          vm.GPUs,
			CPU:           vm.CPU,
			Memory:        vm.Memory,
			RootfsSize:    vm.RootfsSize,
			Entrypoint:    vm.Entrypoint,
			Mounts:        convertToWorkloadMounts(vm.Mounts),
			Zlogs:         convertToWorkloadZlogs(vm.Zlogs),
			EnvVars:       vm.EnvVars,
		}
	}

	return result
}

func convertVMsComputed(VMs []workloads.VM) []VMComputed {

	result := make([]VMComputed, len(VMs))

	for i, vm := range VMs {

		result[i] = VMComputed{
			ComputedIP:  vm.ComputedIP,
			ComputedIP6: vm.ComputedIP6,
			YggIP:       vm.YggIP,
			ConsoleURL:  vm.ConsoleURL,
			IP:          vm.IP,
		}
	}
	return result
}

func convertBackends(backends workloads.Backends) []Backend {
	result := make([]Backend, len(backends))
	for i, backend := range backends {
		result[i] = Backend{
			Address:   backend.Address,
			Namespace: backend.Namespace,
			Password:  backend.Password,
		}
	}
	return result
}

func convertToWorkloadBackends(backends []Backend) workloads.Backends {
	result := make(workloads.Backends, len(backends))
	for i, backend := range backends {
		result[i] = workloads.Backend{
			Address:   backend.Address,
			Namespace: backend.Namespace,
			Password:  backend.Password,
		}
	}
	return result
}

func convertMetadata(metadata workloads.Metadata) Metadata {
	return Metadata{
		EncryptionKey:       metadata.EncryptionKey,
		Prefix:              metadata.Prefix,
		EncryptionAlgorithm: metadata.EncryptionAlgorithm,
		Type:                metadata.Type,
		Backends:            convertBackends(metadata.Backends),
	}
}

func convertToWorkloadMetadata(metadata Metadata) workloads.Metadata {
	return workloads.Metadata{
		EncryptionKey:       metadata.EncryptionKey,
		Prefix:              metadata.Prefix,
		EncryptionAlgorithm: metadata.EncryptionAlgorithm,
		Type:                metadata.Type,
		Backends:            convertToWorkloadBackends(metadata.Backends),
	}
}

func convertGroups(groups workloads.Groups) []Group {
	result := make([]Group, len(groups))
	for i, group := range groups {
		result[i] = Group{
			Backends: convertBackends(group.Backends),
		}
	}
	return result
}

func convertToWorkloadGroups(groups []Group) workloads.Groups {
	result := make(workloads.Groups, len(groups))
	for i, group := range groups {
		result[i] = workloads.Group{
			Backends: convertToWorkloadBackends(group.Backends),
		}
	}
	return result
}

func convertQSFSs(qsfss []workloads.QSFS) []QSFSInput {

	result := make([]QSFSInput, len(qsfss))

	for i, qsfs := range qsfss {
		result[i] = QSFSInput{
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
			Metadata:             convertMetadata(qsfs.Metadata),
			Groups:               convertGroups(qsfs.Groups),
		}
	}
	return result
}

func convertToWorkloadQSFSs(qsfss []QSFSInput) []workloads.QSFS {

	result := make([]workloads.QSFS, len(qsfss))

	for i, qsfs := range qsfss {
		result[i] = workloads.QSFS{
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
			Metadata:             convertToWorkloadMetadata(qsfs.Metadata),
			Groups:               convertToWorkloadGroups(qsfs.Groups),
		}
	}
	return result
}

func convertQSFSsComputed(qsfss []workloads.QSFS) []QSFSComputed {

	result := make([]QSFSComputed, len(qsfss))

	for i, qsfs := range qsfss {
		result[i] = QSFSComputed{
			MetricsEndpoint: qsfs.MetricsEndpoint,
		}
	}
	return result
}

func parseToDeploymentState(deployment workloads.Deployment) DeploymentState {

	stateArgs := DeploymentArgs{

		NodeID:           int32(deployment.NodeID),
		Name:             deployment.Name,
		SolutionType:     deployment.SolutionType,
		SolutionProvider: deployment.SolutionProvider,
		NetworkName:      deployment.NetworkName,
		Disks:            convertDisks(deployment.Disks),
		ZdbsInputs:       convertZdbs(deployment.Zdbs),
		VmsInputs:        convertVMs(deployment.Vms),
		QSFSInputs:       convertQSFSs(deployment.QSFS),
	}

	state := DeploymentState{

		DeploymentArgs:   stateArgs,
		NodeDeploymentID: deployment.NodeDeploymentID,
		ContractID:       deployment.ContractID,
		IPrange:          deployment.IPrange,
		ZdbsComputed:     convertZdbsComputed(deployment.Zdbs),
		VmsComputed:      convertVMsComputed(deployment.Vms),
		QsfsComputed:     convertQSFSsComputed(deployment.QSFS),
	}

	return state
}

func parseToWorkloadDeployment(deploymentArgs DeploymentArgs) workloads.Deployment {

	return workloads.Deployment{
		NodeID:           uint32(deploymentArgs.NodeID),
		Name:             deploymentArgs.Name,
		SolutionType:     deploymentArgs.SolutionType,
		SolutionProvider: deploymentArgs.SolutionProvider,
		NetworkName:      deploymentArgs.NetworkName,
		Disks:            convertToWorkloadDisks(deploymentArgs.Disks),
		Zdbs:             convertToWorkloadZdbs(deploymentArgs.ZdbsInputs),
		Vms:              convertToWorkloadVMs(deploymentArgs.VmsInputs),
		QSFS:             convertToWorkloadQSFSs(deploymentArgs.QSFSInputs),
	}

}

// Create creates a deployment
func (*Deployment) Create(ctx p.Context, name string, input DeploymentArgs, preview bool) (string, DeploymentState, error) {

	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return name, state, nil
	}

	deployment := parseToWorkloadDeployment(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return name, state, err
	}

	state = parseToDeploymentState(deployment)

	return name, state, nil
}

// Update updates the arguments of the deployment resource
func (*Deployment) Update(ctx p.Context, name string, input DeploymentArgs, preview bool) (string, DeploymentState, error) {

	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return name, state, nil
	}

	deployment := parseToWorkloadDeployment(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return name, state, err
	}

	state = parseToDeploymentState(deployment)

	return name, state, nil

}

// Read gets the state of the deployment resource
func (*Deployment) Read(ctx p.Context, name string, input DeploymentArgs, preview bool) (string, DeploymentState, error) {

	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return name, state, nil
	}

	deployment := parseToWorkloadDeployment(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return name, state, err
	}

	state = parseToDeploymentState(deployment)

	return name, state, nil

}

// Delete deletes a deployment resource
func (*Deployment) Delete(ctx p.Context, name string, input DeploymentArgs, preview bool) (string, DeploymentState, error) {

	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return name, state, nil
	}

	deployment := parseToWorkloadDeployment(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Cancel(ctx, &deployment); err != nil {
		state = parseToDeploymentState(deployment)
		return name, state, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		state = parseToDeploymentState(deployment)
		return name, state, err
	}

	return name, state, nil
}
