package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	DefaultLabelWidth = 10
	DefaultBarWidth   = 30
	thresholdWarning  = 60.0
	thresholdCritical = 85.0

	barFilledChar = "█"
	barEmptyChar  = "░"
)

func ColorForPercent(percent float64) lipgloss.Color {
	switch {
	case percent >= thresholdCritical:
		return ColorError
	case percent >= thresholdWarning:
		return ColorWarning
	default:
		return ColorSuccess
	}
}

func styleForPercent(percent float64) lipgloss.Style {
	switch {
	case percent >= thresholdCritical:
		return BarFilledCritical
	case percent >= thresholdWarning:
		return BarFilledWarning
	default:
		return BarFilledSuccess
	}
}

func RenderProgressBar(percent float64) string {
	return RenderBar(percent, DefaultBarWidth)
}

func RenderBar(percent float64, width int) string {
	if width <= 0 {
		width = DefaultBarWidth
	}

	p := percent
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}

	filled := int(p / 100 * float64(width))
	if filled > width {
		filled = width
	}

	filledPart := styleForPercent(p).Render(strings.Repeat(barFilledChar, filled))
	emptyPart := BarEmpty.Render(strings.Repeat(barEmptyChar, width-filled))

	percentStyle := PercentStyle.Foreground(ColorForPercent(p))

	return fmt.Sprintf("[%s%s] %s",
		filledPart,
		emptyPart,
		percentStyle.Render(fmt.Sprintf("%5.1f%%", p)),
	)
}

func RenderMetric(label string, percent float64) string {
	return RenderMetricRow(label, percent, 30, 10)
}

func RenderMetricRow(label string, percent float64, width, labelWidth int) string {
	return fmt.Sprintf("%s %s",
		LabelStyle.Width(labelWidth).Render(label),
		RenderBar(percent, width),
	)
}
