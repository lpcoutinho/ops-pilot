package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type GetSystemHealthTool struct{}

type SystemHealthReport struct {
	CPUUsage    float64 `json:"cpu_usage_percent"`
	MemoryTotal uint64  `json:"memory_total_bytes"`
	MemoryUsed  uint64  `json:"memory_used_bytes"`
	MemoryFree  uint64  `json:"memory_free_bytes"`
	DiskUsage   float64 `json:"disk_usage_percent"`
}

func (t *GetSystemHealthTool) Name() string {
	return "get_system_health"
}

func (t *GetSystemHealthTool) Description() string {
	return "Collects current CPU, Memory, and Disk usage metrics from the system."
}

func (t *GetSystemHealthTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *GetSystemHealthTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	report := &SystemHealthReport{}

	// CPU
	cpuPercent, err := cpu.PercentWithContext(ctx, 0, false)
	if err == nil && len(cpuPercent) > 0 {
		report.CPUUsage = cpuPercent[0]
	}

	// Memory
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err == nil {
		report.MemoryTotal = vm.Total
		report.MemoryUsed = vm.Used
		report.MemoryFree = vm.Available
	}

	// Disk
	du, err := disk.UsageWithContext(ctx, "/")
	if err == nil {
		report.DiskUsage = du.UsedPercent
	}

	if err != nil {
		return nil, fmt.Errorf("failed to collect some metrics: %w", err)
	}

	return report, nil
}
