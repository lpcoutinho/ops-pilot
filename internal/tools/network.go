package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
)

type AuditNetworkTool struct{}

type NetworkInterfaceInfo struct {
	Name        string   `json:"name"`
	MTU         int      `json:"mtu"`
	Flags       []string `json:"flags"`
	IPAddresses []string `json:"ip_addresses"`
}

type NetworkStats struct {
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

type NetworkAuditReport struct {
	Interfaces []NetworkInterfaceInfo `json:"interfaces"`
	IOStats    []net.IOCountersStat   `json:"io_stats"`
}

func (t *AuditNetworkTool) Name() string {
	return "audit_network"
}

func (t *AuditNetworkTool) Description() string {
	return "Provides a detailed audit of network interfaces, IP addresses, and traffic statistics (I/O counters)."
}

func (t *AuditNetworkTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *AuditNetworkTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	interfaces, err := net.InterfacesWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	ioStats, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network I/O stats: %w", err)
	}

	var interfaceInfos []NetworkInterfaceInfo
	for _, i := range interfaces {
		var ips []string
		for _, addr := range i.Addrs {
			ips = append(ips, addr.Addr)
		}

		interfaceInfos = append(interfaceInfos, NetworkInterfaceInfo{
			Name:        i.Name,
			MTU:         i.MTU,
			Flags:       i.Flags,
			IPAddresses: ips,
		})
	}

	return NetworkAuditReport{
		Interfaces: interfaceInfos,
		IOStats:    ioStats,
	}, nil
}
