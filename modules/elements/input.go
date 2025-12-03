package elements

import (
	"go-tui-library/modules/core"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type InputElement struct {
	elem   core.Element
	value  string
	cursor int
}

func NewInputElement(elem core.Element) *InputElement {
	return &InputElement{
		elem:  elem,
		value: elem.Value,
	}
}

func (i *InputElement) HandleKey(key string) {
	switch key {
	case "backspace":
		if i.cursor > 0 && len(i.value) > 0 {
			i.value = i.value[:i.cursor-1] + i.value[i.cursor:]
			i.cursor--
		}
	case "left":
		if i.cursor > 0 {
			i.cursor--
		}
	case "right":
		if i.cursor < len(i.value) {
			i.cursor++
		}
	default:
		if len(key) == 1 || key == " " {
			i.value = i.value[:i.cursor] + key + i.value[i.cursor:]
			i.cursor++
		}
	}
}
func (i *InputElement) Render(focused bool) string {
	borderStyle := core.ParseBorderStyle(i.elem.Style.Border)
	if i.elem.Style.Border == "" {
		borderStyle = lipgloss.RoundedBorder()
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		Padding(0, 1).
		Width(30)

	if i.elem.Style.Color != "" {
		style = style.Foreground(core.ParseColor(i.elem.Style.Color))
	}
	if i.elem.Style.Background != "" {
		style = style.Background(core.ParseColor(i.elem.Style.Background))
	}

	if focused {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	displayValue := i.value
	if focused && i.cursor <= len(i.value) {
		if i.cursor < len(i.value) {
			displayValue = i.value[:i.cursor] + "│" + i.value[i.cursor:]
		} else {
			displayValue = i.value + "│"
		}
	}

	if i.elem.Label != "" {
		labelWidth := len(i.elem.Label) + 2
		leftPadding := 2
		rightPadding := style.GetWidth() - labelWidth - leftPadding
		if rightPadding < 0 {
			rightPadding = 0
		}

		topBorder := borderStyle.Top
		if topBorder == "" {
			topBorder = "─"
		}

		topLine := borderStyle.TopLeft +
			strings.Repeat(topBorder, leftPadding) +
			" " + i.elem.Label + " " +
			strings.Repeat(topBorder, rightPadding) +
			borderStyle.TopRight

		middleLine := borderStyle.Left +
			" " + displayValue +
			strings.Repeat(" ", style.GetWidth()-len(displayValue)-2) +
			" " + borderStyle.Right

		bottomLine := borderStyle.BottomLeft +
			strings.Repeat(borderStyle.Bottom, style.GetWidth()) +
			borderStyle.BottomRight

		if focused {
			borderColor := lipgloss.Color("205")
			coloredStyle := lipgloss.NewStyle().Foreground(borderColor)
			topLine = coloredStyle.Render(topLine)
			middleLine = coloredStyle.Render(middleLine)
			bottomLine = coloredStyle.Render(bottomLine)
		}

		return topLine + "\n" + middleLine + "\n" + bottomLine
	}

	return style.Render(displayValue)
}
