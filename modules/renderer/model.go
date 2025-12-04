package renderer

import (
	"strings"

	"go-tui-library/modules/core"
	"go-tui-library/modules/elements"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	config       core.UIConfig
	width        int
	height       int
	focusIndex   int
	focusableIDs []string
	inputs       map[string]*elements.InputElement
	checkboxes   map[string]*elements.CheckboxElement
}

func NewModel(config core.UIConfig) Model {
	m := Model{
		config:     config,
		inputs:     make(map[string]*elements.InputElement),
		checkboxes: make(map[string]*elements.CheckboxElement),
	}
	m.collectFocusableElements(config.Elements)
	return m
}

func (m *Model) collectFocusableElements(elems []core.Element) {
	for i := range elems {
		elem := &elems[i]
		if elem.Type == "input" {
			m.focusableIDs = append(m.focusableIDs, elem.ID)
			m.inputs[elem.ID] = elements.NewInputElement(*elem)
		} else if elem.Type == "checkbox" {
			m.focusableIDs = append(m.focusableIDs, elem.ID)
			m.checkboxes[elem.ID] = elements.NewCheckboxElement(*elem)
		} else if elem.Type == "button" {
			m.focusableIDs = append(m.focusableIDs, elem.ID)
		}
		if len(elem.Children) > 0 {
			m.collectFocusableElements(elem.Children)
		}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % len(m.focusableIDs)
		case "shift+tab":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = len(m.focusableIDs) - 1
			}
		case "enter", " ":
			if m.focusIndex < len(m.focusableIDs) {
				focusedID := m.focusableIDs[m.focusIndex]
				if cb, ok := m.checkboxes[focusedID]; ok {
					cb.Toggle()
				}
			}
		default:
			if m.focusIndex < len(m.focusableIDs) {
				focusedID := m.focusableIDs[m.focusIndex]
				if input, ok := m.inputs[focusedID]; ok {
					input.HandleKey(msg.String())
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	layout := m.config.Layout
	if layout == "" {
		layout = "column"
	}

	var content string
	if layout == "flex" {
		content = m.renderFlexLayout(m.config.Elements)
	} else {
		content = m.renderColumnLayout(m.config.Elements)
	}

	return content
}

func (m Model) renderColumnLayout(elems []core.Element) string {
	var columns []string
	for _, elem := range elems {
		rendered := m.renderElement(elem, m.width)
		columns = append(columns, rendered)
	}
	return lipgloss.JoinVertical(lipgloss.Left, columns...)
}

func (m Model) renderFlexLayout(elems []core.Element) string {
	var items []string
	for _, elem := range elems {
		width := m.width / len(elems)
		if elem.Width > 0 {
			width = elem.Width
		}
		rendered := m.renderElement(elem, width)
		items = append(items, rendered)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, items...)
}

func (m Model) renderElement(elem core.Element, width int) string {
	focused := false
	if m.focusIndex < len(m.focusableIDs) {
		focused = m.focusableIDs[m.focusIndex] == elem.ID
	}

	switch elem.Type {
	case "div":
		return m.renderDiv(elem, width)
	case "text":
		return m.renderText(elem)
	case "input":
		if input, ok := m.inputs[elem.ID]; ok {
			return input.Render(focused)
		}
		return ""
	case "checkbox":
		if cb, ok := m.checkboxes[elem.ID]; ok {
			return cb.Render(focused, width)
		}
		return ""
	case "button":
		return m.renderButton(elem, focused)
	case "table":
		return m.renderTable(elem)
	default:
		return ""
	}
}

func (m Model) renderDiv(elem core.Element, parentWidth int) string {
	layout := elem.Layout
	if layout == "" {
		layout = "column"
	}

	// Use explicit width if provided, otherwise use parent width
	width := parentWidth
	if elem.Width > 0 {
		width = elem.Width
	}

	// Use explicit height if provided
	height := 0
	if elem.Height > 0 {
		height = elem.Height
	}

	var content string
	if layout == "flex" {
		flexDir := elem.FlexDir
		if flexDir == "" {
			flexDir = "row"
		}
		content = m.renderFlexDiv(elem, width, flexDir)
	} else {
		var children []string
		for _, child := range elem.Children {
			children = append(children, m.renderElement(child, width))
		}
		content = strings.Join(children, "\n")
	}

	style := lipgloss.NewStyle()

	// Apply width
	if width > 0 {
		style = style.Width(width)
	}

	// Apply height
	if height > 0 {
		style = style.Height(height)
	}

	if elem.Align != "" {
		style = m.applyAlign(style, elem.Align)
	}

	if elem.Style.Color != "" {
		style = style.Foreground(core.ParseColor(elem.Style.Color))
	}
	if elem.Style.Background != "" {
		style = style.Background(core.ParseColor(elem.Style.Background))
	}
	if elem.Style.Border != "" {
		style = style.Border(core.ParseBorderStyle(elem.Style.Border))
	}

	return style.Render(content)
}

func (m Model) renderFlexDiv(elem core.Element, width int, flexDir string) string {
	var items []string

	for _, child := range elem.Children {
		childWidth := width
		if flexDir == "row" && len(elem.Children) > 0 {
			// Distribute width among children
			childWidth = width / len(elem.Children)
		}
		// Override with explicit child width if provided
		if child.Width > 0 {
			childWidth = child.Width
		}
		rendered := m.renderElement(child, childWidth)
		items = append(items, rendered)
	}

	var content string
	if flexDir == "column" {
		align := lipgloss.Left
		if elem.Align == "center" {
			align = lipgloss.Center
		} else if elem.Align == "right" || elem.Align == "flex-end" {
			align = lipgloss.Right
		}
		content = lipgloss.JoinVertical(align, items...)
	} else {
		align := lipgloss.Top
		if elem.Align == "center" {
			align = lipgloss.Center
		} else if elem.Align == "bottom" || elem.Align == "flex-end" {
			align = lipgloss.Bottom
		}
		content = lipgloss.JoinHorizontal(align, items...)
	}

	if elem.Justify != "" {
		content = m.applyJustify(content, elem.Justify, width, flexDir, elem.Height)
	}

	return content
}

func (m Model) applyAlign(style lipgloss.Style, align string) lipgloss.Style {
	switch align {
	case "center":
		return style.Align(lipgloss.Center)
	case "right", "flex-end":
		return style.Align(lipgloss.Right)
	case "left", "flex-start":
		return style.Align(lipgloss.Left)
	default:
		return style
	}
}

func (m Model) applyJustify(content string, justify string, width int, flexDir string, elemHeight int) string {
	contentWidth := lipgloss.Width(content)
	contentHeight := lipgloss.Height(content)

	// Use element height if provided, otherwise use screen height
	targetHeight := m.height
	if elemHeight > 0 {
		targetHeight = elemHeight
	}

	if flexDir == "row" {
		switch justify {
		case "center":
			padding := (width - contentWidth) / 2
			if padding > 0 {
				return lipgloss.NewStyle().PaddingLeft(padding).Render(content)
			}
		case "flex-end", "right":
			padding := width - contentWidth
			if padding > 0 {
				return lipgloss.NewStyle().PaddingLeft(padding).Render(content)
			}
		case "space-between":
			return lipgloss.NewStyle().Width(width).Render(content)
		case "space-around":
			padding := (width - contentWidth) / 2
			if padding > 0 {
				return lipgloss.NewStyle().PaddingLeft(padding).PaddingRight(padding).Render(content)
			}
		}
	} else {
		switch justify {
		case "center":
			padding := (targetHeight - contentHeight) / 2
			if padding > 0 {
				return lipgloss.NewStyle().PaddingTop(padding).Render(content)
			}
		case "flex-end", "bottom":
			padding := targetHeight - contentHeight
			if padding > 0 {
				return lipgloss.NewStyle().PaddingTop(padding).Render(content)
			}
		}
	}

	return content
}

func (m Model) renderText(elem core.Element) string {
	style := lipgloss.NewStyle()

	// Apply width if specified
	if elem.Width > 0 {
		style = style.Width(elem.Width)
	}

	// Apply height if specified
	if elem.Height > 0 {
		style = style.Height(elem.Height)
	}

	if elem.Style.Color != "" {
		style = style.Foreground(core.ParseColor(elem.Style.Color))
	}
	if elem.Style.Background != "" {
		style = style.Background(core.ParseColor(elem.Style.Background))
	}
	return style.Render(elem.Value)
}

func (m Model) renderButton(elem core.Element, focused bool) string {
	style := lipgloss.NewStyle().
		Padding(0, 2).
		Border(core.ParseBorderStyle(elem.Style.Border))

	// Apply width if specified
	if elem.Width > 0 {
		style = style.Width(elem.Width)
	}

	// Apply height if specified
	if elem.Height > 0 {
		style = style.Height(elem.Height)
	}

	if elem.Style.Color != "" {
		style = style.Foreground(core.ParseColor(elem.Style.Color))
	}
	if elem.Style.Background != "" {
		style = style.Background(core.ParseColor(elem.Style.Background))
	}

	if focused {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	label := elem.Label
	if label == "" {
		label = "Button"
	}

	return style.Render(label)
}

func (m Model) renderTable(elem core.Element) string {
	if len(elem.Headers) == 0 {
		return ""
	}

	borderStyle := core.ParseBorderStyle(elem.Style.Border)
	if elem.Style.Border == "" {
		borderStyle = lipgloss.NormalBorder()
	}

	headerStyle := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(core.ParseColor(elem.Style.Color)).
		Padding(0, 1).
		Bold(true)

	if elem.Style.Color != "" {
		headerStyle = headerStyle.Foreground(core.ParseColor(elem.Style.Color))
	}
	if elem.Style.Background != "" {
		headerStyle = headerStyle.Background(core.ParseColor(elem.Style.Background))
	}

	cellStyle := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(core.ParseColor(elem.Style.Color)).
		Padding(0, 1)

	if elem.Style.Color != "" {
		cellStyle = cellStyle.Foreground(core.ParseColor(elem.Style.Color))
	}

	var headerRow []string
	for _, header := range elem.Headers {
		headerRow = append(headerRow, headerStyle.Render(header))
	}

	var rows []string
	rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, headerRow...))

	for _, row := range elem.Rows {
		var cells []string
		for i := 0; i < len(elem.Headers); i++ {
			cellValue := ""
			if i < len(row) {
				cellValue = row[i]
			}
			cells = append(cells, cellStyle.Render(cellValue))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
