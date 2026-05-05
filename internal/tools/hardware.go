package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

type GetHardwareInfoTool struct{}

type HardwareReport struct {
	Hostname        string `json:"hostname"`
	Uptime          uint64 `json:"uptime_seconds"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	KernelVersion   string `json:"kernel_version"`
	KernelArch      string `json:"kernel_arch"`
	CPUModel        string `json:"cpu_model"`
	CPUCoresPhysical int    `json:"cpu_cores_physical"`
	CPUCoresLogical  int    `json:"cpu_cores_logical"`
}

func (t *GetHardwareInfoTool) Name() string {
	return "get_hardware_info"
}

func (t *GetHardwareInfoTool) Description() string {
	return "Returns detailed information about the system hardware, including CPU model, OS version, kernel, and uptime."
}

func (t *GetHardwareInfoTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *GetHardwareInfoTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	hInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}

	cInfo, err := cpu.InfoWithContext(ctx)
	var cpuModel string
	if err == nil && len(cInfo) > 0 {
		cpuModel = cInfo[0].ModelName
	}

	physicalCores, _ := cpu.CountsWithContext(ctx, false)
	logicalCores, _ := cpu.CountsWithContext(ctx, true)

	return HardwareReport{
		Hostname:         hInfo.Hostname,
		Uptime:           hInfo.Uptime,
		OS:               hInfo.OS,
		Platform:         hInfo.Platform,
		KernelVersion:    hInfo.KernelVersion,
		KernelArch:       hInfo.KernelArch,
		CPUModel:         cpuModel,
		CPUCoresPhysical: physicalCores,
		CPUCoresLogical:  logicalCores,
	}, nil
}
