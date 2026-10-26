package metrics

import (
	"runtime"

	"github.com/Iposhka54/pulse/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
)

var buildInfo = prometheus.NewDesc(
	"pulse_build_info",
	"Build information about the pulse binary.",
	[]string{"goversion"}, nil,
)

var cpuCores = prometheus.NewDesc(
	"pulse_cpu_cores",
	"Number of CPU cores.",
	nil, nil,
)

var (
	cpuUsageRatio = prometheus.NewDesc(
		"pulse_cpu_usage_ratio",
		"Current CPU usage ratio (0-1).",
		nil, nil,
	)
	memUsed = prometheus.NewDesc(
		"pulse_memory_used_bytes",
		"Memory used in bytes.",
		nil, nil,
	)
	memTotal = prometheus.NewDesc(
		"pulse_memory_total_bytes",
		"Total memory in bytes.",
		nil, nil,
	)
	diskUsed = prometheus.NewDesc(
		"pulse_disk_used_bytes",
		"Disk used in bytes.",
		[]string{"mountpoint"}, nil,
	)
	diskTotal = prometheus.NewDesc(
		"pulse_disk_total_bytes",
		"Disk total in bytes.",
		[]string{"mountpoint"}, nil,
	)
	load1 = prometheus.NewDesc(
		"pulse_load1",
		"1-minute load average.",
		nil, nil,
	)
	load5 = prometheus.NewDesc(
		"pulse_load5",
		"5-minute load average.",
		nil, nil,
	)
	load15 = prometheus.NewDesc(
		"pulse_load15",
		"15-minute load average.",
		nil, nil,
	)
	uptimeSeconds = prometheus.NewDesc(
		"pulse_uptime_seconds",
		"System uptime in seconds.",
		nil, nil,
	)
)

type SystemCollector struct {
	src *collector.Collector
}

func NewSystemCollector(src *collector.Collector) *SystemCollector {
	return &SystemCollector{src: src}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- buildInfo
	ch <- cpuUsageRatio
	ch <- cpuCores
	ch <- memUsed
	ch <- memTotal
	ch <- diskUsed
	ch <- diskTotal
	ch <- load1
	ch <- load5
	ch <- load15
	ch <- uptimeSeconds
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		buildInfo, prometheus.GaugeValue, 1, runtime.Version(),
	)
	snap := c.src.Latest()
	if snap == nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(cpuUsageRatio, prometheus.GaugeValue, snap.CPU.Percent/100)
	ch <- prometheus.MustNewConstMetric(cpuCores, prometheus.GaugeValue, float64(snap.CPU.Cores))
	ch <- prometheus.MustNewConstMetric(memUsed, prometheus.GaugeValue, float64(snap.Memory.Used))
	ch <- prometheus.MustNewConstMetric(memTotal, prometheus.GaugeValue, float64(snap.Memory.Total))

	for _, d := range snap.Disks {
		ch <- prometheus.MustNewConstMetric(diskUsed, prometheus.GaugeValue, float64(d.Used), d.Mountpoint)
		ch <- prometheus.MustNewConstMetric(diskTotal, prometheus.GaugeValue, float64(d.Total), d.Mountpoint)
	}

	if snap.Load.L1 > 0 || snap.Load.L5 > 0 || snap.Load.L15 > 0 {
		ch <- prometheus.MustNewConstMetric(load1, prometheus.GaugeValue, snap.Load.L1)
		ch <- prometheus.MustNewConstMetric(load5, prometheus.GaugeValue, snap.Load.L5)
		ch <- prometheus.MustNewConstMetric(load15, prometheus.GaugeValue, snap.Load.L15)
	}

	ch <- prometheus.MustNewConstMetric(uptimeSeconds, prometheus.GaugeValue, snap.Uptime.Seconds())
}
