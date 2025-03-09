package monitor

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type (
	// SystemInfo 包含所有系统性能信息
	SystemInfo struct {
		CPUUsagePercent []float64              `json:"cpuUsagePercent"`
		Memory          *mem.VirtualMemoryStat `json:"memory"`
		DiskUsage       *disk.UsageStat        `json:"diskUsage"`
		NetStats        []net.IOCountersStat   `json:"netStats"`
		LoadAvg         *load.AvgStat          `json:"loaAvg"`
		SysInfo         *host.InfoStat         `json:"sysInfo"`
	}

	// NetworkInfo 网络接口流量
	NetworkInfo struct {
		Name        string `json:"name"`
		BytesSent   uint64 `json:"bytesSent"`
		BytesRecv   uint64 `json:"bytesRecv"`
		PacketsSent uint64 `json:"packetsSent"`
		PacketsRecv uint64 `json:"packetsRecv"`
	}

	IMonitorService interface {
		GetStatus() (resp *SystemInfo, err error)
	}
	monitorServiceImpl struct{}
)

var (
	_ IMonitorService = (*monitorServiceImpl)(nil)

	DefaultMonitorService IMonitorService = &monitorServiceImpl{}
)

func (*monitorServiceImpl) GetStatus() (resp *SystemInfo, err error) {
	// 构建 SystemInfo 结构体
	resp = &SystemInfo{}

	if cpuNum, e := cpu.Percent(0, true); e == nil {
		resp.CPUUsagePercent = cpuNum
	}

	// 获取内存使用情况
	if vMem, e := mem.VirtualMemory(); e == nil {
		resp.Memory = vMem
	}

	// 获取磁盘使用情况
	if diskUsage, e := disk.Usage("/"); e == nil {
		resp.DiskUsage = diskUsage
	}

	// 获取系统负载
	if loadAvg, e := load.Avg(); e == nil {
		resp.LoadAvg = loadAvg
	}

	// 获取系统信息
	if sysInfo, e := host.Info(); e == nil {
		resp.SysInfo = sysInfo
	}

	// 获取网络接口流量
	if netStats, e := net.IOCounters(true); e == nil {
		resp.NetStats = netStats
	}

	return
}
