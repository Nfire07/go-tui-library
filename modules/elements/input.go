package elements

import (
	"go-tui-library/modules/core"

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
		if len(key) == 1 {
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

	if i.elem.Label != "" {
		style = style.BorderTop(true).BorderTopForeground(lipgloss.Color("240"))
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
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		return labelStyle.Render(i.elem.Label) + "\n" + style.Render(displayValue)
	}

	return style.Render(displayValue)
}
