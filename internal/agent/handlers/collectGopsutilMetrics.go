package handlers

import (
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// CollectGopsutilMetrics collects various metrics to given storage
func CollectGopsutilMetrics(s Storage, collectInterval time.Duration) error {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	s.UpdateGauge("TotalMemory", float64(vmem.Total))
	s.UpdateGauge("FreeMemory", float64(vmem.Free))
	cpus, err := cpu.Percent(collectInterval, true)
	if err != nil {
		return err
	}
	for i, percent := range cpus {
		s.UpdateGauge("CPUutilization"+strconv.Itoa(i+1), percent)
	}
	return nil
}
