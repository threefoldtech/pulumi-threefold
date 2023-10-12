package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

// Scheduler controlling struct
type Scheduler struct{}

// SchedulerArgs is defining what arguments it accepts
type SchedulerArgs struct {
	MRU               int64   `pulumi:"mru,optional"`
	HRU               int64   `pulumi:"hru,optional"`
	SRU               int64   `pulumi:"sru,optional"`
	Country           string  `pulumi:"country,optional"`
	City              string  `pulumi:"city,optional"`
	FarmName          string  `pulumi:"farm_name,optional"`
	FarmIDs           []int64 `pulumi:"farm_ids,optional"`
	FreeIPs           int64   `pulumi:"free_ips,optional"`
	IPv4              bool    `pulumi:"ipv4,optional"`
	IPv6              bool    `pulumi:"ipv6,optional"`
	Domain            bool    `pulumi:"domain,optional"`
	Dedicated         bool    `pulumi:"dedicated,optional"`
	Rented            bool    `pulumi:"rented,optional"`
	Rentable          bool    `pulumi:"rentable,optional"`
	NodeID            int64   `pulumi:"node_id,optional"`
	TwinID            int64   `pulumi:"twin_id,optional"`
	CertificationType string  `pulumi:"certification_type,optional"`
	HasGPU            bool    `pulumi:"has_gpu,optional"`
	GpuDeviceID       string  `pulumi:"gpu_device_id,optional"`
	GpuDeviceName     string  `pulumi:"gpu_device_name,optional"`
	GpuVendorID       string  `pulumi:"gpu_vendor_id,optional"`
	GpuVendorName     string  `pulumi:"gpu_vendor_name,optional"`
	GpuAvailable      bool    `pulumi:"gpu_available,optional"`
}

// SchedulerState is describing the fields that exist on the created resource.
type SchedulerState struct {
	SchedulerArgs

	Nodes []int32 `pulumi:"nodes"`
}

// Create creates scheduler
func (*Scheduler) Create(ctx p.Context, id string, input SchedulerArgs, preview bool) (string, SchedulerState, error) {
	state := SchedulerState{SchedulerArgs: input}
	if preview {
		return id, state, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodeFilter, ssds, hdds := parseSchedulerInput(input)

	nodes, err := deployer.FilterNodes(ctx, config.TFPluginClient, nodeFilter, hdds, ssds, nil)
	if err != nil {
		return id, state, err
	}

	for _, node := range nodes {
		state.Nodes = append(state.Nodes, int32(node.NodeID))
	}

	return id, state, nil
}

// Update updates the arguments of the scheduler resource
func (*Scheduler) Update(
	ctx p.Context,
	id string,
	oldState SchedulerState,
	input SchedulerArgs,
	preview bool) (SchedulerState, error) {
	state := SchedulerState{SchedulerArgs: input}
	if preview {
		return state, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodeFilter, hdds, ssds := parseSchedulerInput(input)
	nodes, err := deployer.FilterNodes(ctx, config.TFPluginClient, nodeFilter, hdds, ssds, nil)
	if err != nil {
		return state, err
	}

	for _, node := range nodes {
		state.Nodes = append(state.Nodes, int32(node.NodeID))
	}

	return state, nil
}

// Read get the state of the scheduler resource
func (*Scheduler) Read(ctx p.Context, id string, oldState SchedulerState) (string, SchedulerState, error) {
	return id, oldState, nil
}

// Delete deletes the scheduler resource
func (*Scheduler) Delete(ctx p.Context, id string, oldState SchedulerState) error {
	return nil
}
