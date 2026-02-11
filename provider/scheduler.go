package provider

import (
	"context"
	"errors"
	"slices"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/zos"
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
	Ygg               bool    `pulumi:"ygg,optional"`
	Wireguard         bool    `pulumi:"wireguard,optional"`
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
func (*Scheduler) Create(
	ctx context.Context,
	req infer.CreateRequest[SchedulerArgs],
) (infer.CreateResponse[SchedulerState], error) {
	state := SchedulerState{SchedulerArgs: req.Inputs}
	if req.DryRun {
		return infer.CreateResponse[SchedulerState]{ID: req.Name, Output: state}, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodeFilter, ssds, hdds := parseSchedulerInput(req.Inputs)

	nodes, err := deployer.FilterNodes(ctx, config.TFPluginClient, nodeFilter, hdds, ssds, nil)
	if errors.Is(err, deployer.ErrNoNodesMatchesResources) && slices.Contains(nodeFilter.Features, zos.NetworkLightType) {
		nodeFilter.Features = []string{zos.NetworkType, zos.ZMachineType}
		nodes, err = deployer.FilterNodes(ctx, config.TFPluginClient, nodeFilter, hdds, ssds, nil)
	}

	if err != nil {
		return infer.CreateResponse[SchedulerState]{ID: req.Name, Output: state}, err
	}

	for _, node := range nodes {
		state.Nodes = append(state.Nodes, int32(node.NodeID))
	}

	return infer.CreateResponse[SchedulerState]{ID: req.Name, Output: state}, nil
}

// Update updates the arguments of the scheduler resource
func (*Scheduler) Update(
	ctx context.Context,
	req infer.UpdateRequest[SchedulerArgs, SchedulerState],
) (infer.UpdateResponse[SchedulerState], error) {
	state := SchedulerState{SchedulerArgs: req.Inputs}
	if req.DryRun {
		return infer.UpdateResponse[SchedulerState]{Output: state}, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodeFilter, hdds, ssds := parseSchedulerInput(req.Inputs)
	nodes, err := deployer.FilterNodes(ctx, config.TFPluginClient, nodeFilter, hdds, ssds, nil)
	if err != nil {
		return infer.UpdateResponse[SchedulerState]{Output: state}, err
	}

	for _, node := range nodes {
		state.Nodes = append(state.Nodes, int32(node.NodeID))
	}

	return infer.UpdateResponse[SchedulerState]{Output: state}, nil
}

// Read get the state of the scheduler resource
func (*Scheduler) Read(ctx context.Context, req infer.ReadRequest[SchedulerArgs, SchedulerState]) (infer.ReadResponse[SchedulerArgs, SchedulerState], error) {
	return infer.ReadResponse[SchedulerArgs, SchedulerState](req), nil
}

// Delete deletes the scheduler resource
func (*Scheduler) Delete(ctx context.Context, req infer.DeleteRequest[SchedulerState]) error {
	return nil
}
