package ui

import "github.com/charmbracelet/lipgloss"

// Theme colors - Modern blue theme inspired by Claude Code
var (
	PrimaryColor = lipgloss.Color("#5B9BD5") //Blue

	// Status colors
	ColorSuccess = lipgloss.Color("#2ECC71") // green
	ColorWarning = lipgloss.Color("#F1C40F") // yellow
	ColorError   = lipgloss.Color("#E74C3C") // red
	ColorMuted   = lipgloss.Color("#3A3A3A") // dark-gray
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

	BarFilledSuccess = lipgloss.NewStyle().
				Foreground(ColorSuccess)

	BarFilledWarning = lipgloss.NewStyle().
				Foreground(ColorWarning)

	BarFilledCritical = lipgloss.NewStyle().
				Foreground(ColorError)

	BarEmpty = lipgloss.NewStyle().
			Foreground(ColorMuted)

	PercentStyle = lipgloss.NewStyle().
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor)
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
