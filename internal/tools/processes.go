package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

type GetTopProcessesTool struct{}

type ProcessInfo struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float32 `json:"memory_percent"`
}

type TopProcessesReport struct {
	TopCPU    []ProcessInfo `json:"top_cpu"`
	TopMemory []ProcessInfo `json:"top_memory"`
}

func (t *GetTopProcessesTool) Name() string {
	return "get_top_processes"
}

func (t *GetTopProcessesTool) Description() string {
	return "Returns a list of the top 5 processes consuming the most CPU and Memory currently running on the system."
}

func (t *GetTopProcessesTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *GetTopProcessesTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	procs, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var procInfos []ProcessInfo
	for _, p := range procs {
		name, _ := p.NameWithContext(ctx)
		cpuPerc, _ := p.CPUPercentWithContext(ctx)
		memPerc, _ := p.MemoryPercentWithContext(ctx)

		if name == "" {
			name = "unknown"
		}

		procInfos = append(procInfos, ProcessInfo{
			PID:           p.Pid,
			Name:          name,
			CPUPercent:    cpuPerc,
			MemoryPercent: memPerc,
		})
	}

	if len(procInfos) == 0 {
		return nil, fmt.Errorf("no processes could be read")
	}

	// Sort and pick top 5 for CPU
	sort.Slice(procInfos, func(i, j int) bool {
		return procInfos[i].CPUPercent > procInfos[j].CPUPercent
	})
	
	limit := 5
	if len(procInfos) < limit {
		limit = len(procInfos)
	}
	topCPU := make([]ProcessInfo, limit)
	copy(topCPU, procInfos[:limit])

	// Sort and pick top 5 for Memory
	sort.Slice(procInfos, func(i, j int) bool {
		return procInfos[i].MemoryPercent > procInfos[j].MemoryPercent
	})

	topMemory := make([]ProcessInfo, limit)
	copy(topMemory, procInfos[:limit])

	return TopProcessesReport{
		TopCPU:    topCPU,
		TopMemory: topMemory,
	}, nil
}
