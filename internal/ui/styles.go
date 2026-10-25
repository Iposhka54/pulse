package ui

import "github.com/charmbracelet/lipgloss"

// Theme colors - Modern blue theme inspired by Claude Code
var (
	PrimaryColor = lipgloss.Color("#5B9BD5") //Blue

	// Status colors
	SuccessColor = lipgloss.Color("#2ECC71") // Green
	WarningColor = lipgloss.Color("#F1C40F") // Yellow
	ErrorColor   = lipgloss.Color("#E74C3C") // Red
	InfoColor    = lipgloss.Color("#5B9BD5") // Blue
)

// Base styles
var (
	BoldStyle    = lipgloss.NewStyle().Bold(true)
	PrimaryStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)
)

var (
	BannerStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)
)

func RenderBanner() string {
	banner := `██████╗ ██╗   ██╗██╗     ███████╗███████╗
██╔══██╗██║   ██║██║     ██╔════╝██╔════╝
██████╔╝██║   ██║██║     ███████╗█████╗  
██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝  
██║     ╚██████╔╝███████╗███████║███████╗
╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝`
	return BannerStyle.Render(banner)
}
