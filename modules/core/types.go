package core

type UIConfig struct {
	Layout   string    `json:"layout"`
	Elements []Element `json:"elements"`
}

type Element struct {
	Type     string     `json:"type"`
	ID       string     `json:"id"`
	Label    string     `json:"label,omitempty"`
	Value    string     `json:"value,omitempty"`
	Style    Style      `json:"style,omitempty"`
	Children []Element  `json:"children,omitempty"`
	Width    int        `json:"width,omitempty"`
	Checked  bool       `json:"checked,omitempty"`
	Headers  []string   `json:"headers,omitempty"`
	Rows     [][]string `json:"rows,omitempty"`
}

type Style struct {
	Color      string `json:"color,omitempty"`
	Background string `json:"background,omitempty"`
	Border     string `json:"border,omitempty"`
}

type Position struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}
