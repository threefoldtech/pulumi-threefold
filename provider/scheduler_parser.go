package provider

import (
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/zos"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)

func convertGBToBytes(gb uint64) *uint64 {
	bytes := gb * 1024 * 1024 * 1024
	return &bytes
}

func ref[T any](v T) *T {
	return &v
}

func parseSchedulerInput(input SchedulerArgs) (types.NodeFilter, []uint64, []uint64) {
	filter := types.NodeFilter{Status: []string{"up"}}
	var ssds []uint64
	var hdds []uint64

	for _, farmID := range input.FarmIDs {
		filter.FarmIDs = append(filter.FarmIDs, uint64(farmID))
	}

	if input.MRU != 0 {
		filter.FreeMRU = convertGBToBytes(uint64(input.MRU))
	}

	if input.HRU != 0 {
		filter.FreeHRU = convertGBToBytes(uint64(input.HRU))
		hdds = append(hdds, *filter.FreeHRU)
	}

	if input.SRU != 0 {
		filter.FreeSRU = convertGBToBytes(uint64(input.SRU))
		ssds = append(ssds, *filter.FreeSRU)
	}

	if input.FreeIPs != 0 {
		filter.FreeIPs = ref(uint64(input.FreeIPs))
	}

	if input.NodeID != 0 {
		filter.NodeID = ref(uint64(input.NodeID))
	}

	if input.TwinID != 0 {
		filter.TwinID = ref(uint64(input.TwinID))
	}

	if input.IPv4 {
		filter.IPv4 = &input.IPv4
	}

	if input.IPv6 {
		filter.IPv6 = &input.IPv6
	}

	if input.Domain {
		filter.Domain = &input.Domain
	}

	if input.Dedicated {
		filter.Dedicated = &input.Dedicated
	}

	if input.Rentable {
		filter.Rentable = &input.Rentable
	}

	if input.Rented {
		filter.Rented = &input.Rented
	}

	if input.HasGPU {
		filter.HasGPU = &input.HasGPU
	}

	if input.Rented {
		filter.GpuAvailable = &input.GpuAvailable
	}

	if input.Country != "" {
		filter.Country = &input.Country
	}

	if input.City != "" {
		filter.City = &input.City
	}

	if input.FarmName != "" {
		filter.FarmName = &input.FarmName
	}

	if input.CertificationType != "" {
		filter.CertificationType = &input.CertificationType
	}

	if input.GpuDeviceID != "" {
		filter.GpuDeviceID = &input.GpuDeviceID
	}

	if input.GpuDeviceName != "" {
		filter.GpuDeviceName = &input.GpuDeviceName
	}

	if input.GpuVendorID != "" {
		filter.GpuVendorID = &input.GpuVendorID
	}

	if input.GpuVendorName != "" {
		filter.GpuVendorName = &input.GpuVendorName
	}

	if input.IPv4 || input.IPv6 || input.Domain || input.Ygg || input.Wireguard {
		filter.Features = []string{zos.ZMachineType, zos.NetworkType}
	} else {
		filter.Features = []string{zos.ZMachineLightType, zos.NetworkLightType}
	}

	return filter, ssds, hdds
}
