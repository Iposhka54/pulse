package ui

import (
	"fmt"
	"strings"

	"github.com/Iposhka54/pulse/internal/collector"
)

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func PrintStatus(s *collector.Snapshot) {
	fmt.Println("Host:", s.Hostname)
	fmt.Println("OS:", s.OS)
	fmt.Println("Kernel:", s.Kernel)
	fmt.Println("Uptime:", s.Uptime)
	fmt.Println()

	cpuSubtitle := fmt.Sprintf("%s (%d cores)", s.CPU.Model, s.CPU.Cores)
	fmt.Println(RenderMetricBlock("CPU", cpuSubtitle, s.CPU.Percent))

	memSubtitle := fmt.Sprintf("%s / %s",
		FormatBytes(s.Memory.Used), FormatBytes(s.Memory.Total))
	fmt.Println(RenderMetricBlock("Memory", memSubtitle, s.Memory.Percent))

	if s.Load.L1 > 0 {
		fmt.Printf("%s %.2f %.2f %.2f\n",
			LabelStyle.Width(DefaultLabelWidth).Render("Load"),
			s.Load.L1, s.Load.L5, s.Load.L15,
		)
	}

	fmt.Println()
	fmt.Println("Disks:")
	for _, d := range s.Disks {
		fmt.Printf("  %-20s %s  %s / %s\n",
			d.Mountpoint,
			RenderBar(d.Percent, DefaultBarWidth),
			FormatBytes(d.Used),
			FormatBytes(d.Total),
		)
	}
}

func RenderMetricBlock(label, details string, percent float64) string {
	header := fmt.Sprintf("%s %s",
		LabelStyle.Width(DefaultLabelWidth).Render(label),
		details,
	)
	indent := strings.Repeat(" ", DefaultLabelWidth+1)
	return header + "\n" + indent + RenderBar(percent, DefaultBarWidth)
}
