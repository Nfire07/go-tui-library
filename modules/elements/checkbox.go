package elements

import (
	"go-tui-library/modules/core"

	"github.com/charmbracelet/lipgloss"
)

type CheckboxElement struct {
	elem    core.Element
	checked bool
}

func NewCheckboxElement(elem core.Element) *CheckboxElement {
	return &CheckboxElement{
		elem:    elem,
		checked: elem.Checked,
	}
}

func (c *CheckboxElement) Toggle() {
	c.checked = !c.checked
}

func (c *CheckboxElement) Render(focused bool, width int) string {
	checkbox := "[ ]"
	if c.checked {
		checkbox = "[✓]"
	}

	style := lipgloss.NewStyle()
	if c.elem.Style.Color != "" {
		style = style.Foreground(core.ParseColor(c.elem.Style.Color))
	}
	if c.elem.Style.Background != "" {
		style = style.Background(core.ParseColor(c.elem.Style.Background))
	}

	if focused {
		style = style.Foreground(lipgloss.Color("205"))
	}

	label := c.elem.Label
	if label == "" {
		label = "Checkbox"
	}

	if width > 40 {
		return style.Render(label + " " + checkbox)
	}
	return style.Render(checkbox + "\n" + label)
}
