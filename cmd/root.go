package cmd

import (
	json2 "encoding/json"
	"fmt"
	"os"

	"github.com/Iposhka54/pulse/internal/collector"
	"github.com/Iposhka54/pulse/internal/ui"
	"github.com/spf13/cobra"
)

func Execute() error {
	root := cobra.Command{
		Use:   "pulse",
		Short: "Essential tool for your cluster",
		Long: ui.RenderBanner() + `

Pulse - Professional Server Monitoring Tool

Monitor your servers with real-time metrics and smart notification.
Keep a finger on the pulse of your servers and Kubernetes cluster. One command, instant health.

Need help?
Telegram: @machine_go
Vk: kafka_starter
`,
		DisableSuggestions: true,
		CompletionOptions:  cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	root.AddCommand(newStatusCmd())
	return root.Execute()
}

func newStatusCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current system metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := collector.SnapshotNow()
			if err != nil {
				return err
			}
			if asJSON {
				encoder := json2.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				if err = encoder.Encode(s); err != nil {
					return err
				}
				return nil
			}
			fmt.Println("Host:", s.Hostname)
			fmt.Println("OS:", s.OS)
			fmt.Println("Kernel:", s.Kernel)
			fmt.Println("Uptime:", s.Uptime)
			fmt.Println()

			fmt.Printf("CPU: %.1f%% (%d cores) %s\n", s.CPU.Percent, s.CPU.Cores, s.CPU.Model)
			fmt.Printf("Memory: %s / %s (%.1f%%)\n",
				ui.FormatBytes(s.Memory.Used), ui.FormatBytes(s.Memory.Total), s.Memory.Percent)
			if s.Load.L1 > 0 {
				fmt.Printf("Load: %.2f %.2f %.2f\n", s.Load.L1, s.Load.L5, s.Load.L15)
			}

			fmt.Println()
			fmt.Println("Disks:")
			for _, d := range s.Disks {
				fmt.Printf("  %s: %.1f%% used, %s / %s\n",
					d.Mountpoint, d.Percent, ui.FormatBytes(d.Used), ui.FormatBytes(d.Total))
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(
		&asJSON,
		"json", "j",
		false,
		"Output metrics as JSON",
	)

	return cmd
}
