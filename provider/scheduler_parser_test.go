package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedulerParser(t *testing.T) {
	schedulerInput := SchedulerArgs{
		MRU:               1,
		HRU:               1,
		SRU:               1,
		Country:           "Egypt",
		City:              "Cairo",
		FarmName:          "farm",
		FarmIDs:           []int64{1},
		FreeIPs:           1,
		IPv4:              true,
		IPv6:              true,
		Domain:            true,
		Dedicated:         true,
		Rented:            true,
		Rentable:          true,
		NodeID:            1,
		TwinID:            1,
		CertificationType: "type",
		HasGPU:            true,
		GpuDeviceID:       "1",
		GpuDeviceName:     "gpu",
		GpuVendorID:       "id",
		GpuVendorName:     "name",
		GpuAvailable:      true,
	}

	t.Run("parsing input success", func(t *testing.T) {
		filter, _, _ := parseSchedulerInput(schedulerInput)
		assert.Equal(t, filter.FreeMRU, convertGBToBytes(uint64(schedulerInput.MRU)))
		assert.Equal(t, filter.FreeSRU, convertGBToBytes(uint64(schedulerInput.SRU)))
		assert.Equal(t, filter.FreeHRU, convertGBToBytes(uint64(schedulerInput.HRU)))
		assert.Equal(t, filter.Country, &schedulerInput.Country)
		assert.Equal(t, filter.City, &schedulerInput.City)
	})
}
