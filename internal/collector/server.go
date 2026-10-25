package collector

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type Snapshot struct {
	Hostname string
	OS       string
	Kernel   string
	Uptime   time.Duration
	CPU      struct {
		Percent float64
		Cores   int
		Model   string
	}
	Memory struct {
		Total   uint64
		Used    uint64
		Percent float64
	}
	Disks []DiskInfo
	Load  struct{ L1, L5, L15 float64 }
}

type DiskInfo struct {
	Mountpoint string
	Total      uint64
	Used       uint64
	Percent    float64
}

func SnapshotNow() (*Snapshot, error) {
	s := &Snapshot{}

	if info, err := host.Info(); err == nil {
		s.Hostname = info.Hostname
		s.OS = info.Platform + " " + info.PlatformVersion
		s.Kernel = info.KernelVersion
		s.Uptime = time.Duration(info.Uptime) * time.Second
	}

	if p, err := cpu.Percent(time.Second, false); err == nil && len(p) > 0 {
		s.CPU.Percent = p[0]
	}

	if infos, err := cpu.Info(); err == nil && len(infos) > 0 {
		s.CPU.Model = infos[0].ModelName
	}
	if c, err := cpu.Counts(true); err == nil {
		s.CPU.Cores = c
	}

	if vm, err := mem.VirtualMemory(); err == nil {
		s.Memory.Total = vm.Total
		s.Memory.Used = vm.Used
		s.Memory.Percent = vm.UsedPercent
	}

	if parts, err := disk.Partitions(false); err == nil {
		for _, p := range parts {
			if p.Fstype == "tmpfs" || p.Fstype == "devtmpfs" || p.Fstype == "squashfs" {
				continue
			}
			if u, err := disk.Usage(p.Mountpoint); err == nil {
				s.Disks = append(s.Disks, DiskInfo{
					Mountpoint: p.Mountpoint,
					Total:      u.Total,
					Used:       u.Used,
					Percent:    u.UsedPercent,
				})
			}
		}
	}

	if a, err := load.Avg(); err == nil {
		s.Load.L1, s.Load.L5, s.Load.L15 = a.Load1, a.Load5, a.Load15
	}

	return s, nil
}
