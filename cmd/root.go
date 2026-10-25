package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	root := cobra.Command{
		Use:   "pulse",
		Short: "Essential tool for your cluster",
		Long: `
██████╗ ██╗   ██╗██╗     ███████╗███████╗
██╔══██╗██║   ██║██║     ██╔════╝██╔════╝
██████╔╝██║   ██║██║     ███████╗█████╗  
██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝  
██║     ╚██████╔╝███████╗███████║███████╗
╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝

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
	return root.Execute()
}
