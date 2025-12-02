package core

import "github.com/charmbracelet/lipgloss"

func ParseColor(color string) lipgloss.Color {
	if color == "" {
		return lipgloss.Color("")
	}
	return lipgloss.Color(color)
}

func ParseBorderStyle(border string) lipgloss.Border {
	switch border {
	case "double":
		return lipgloss.DoubleBorder()
	case "rounded":
		return lipgloss.RoundedBorder()
	case "solid":
		return lipgloss.NormalBorder()
	case "thick":
		return lipgloss.ThickBorder()
	default:
		return lipgloss.NormalBorder()
	}
}
