package colors

import (
	"sync"

	termui "github.com/gizak/termui/v3"
)

//Theme uses xterm colors 1-255
type Theme struct {
	base                int
	disabled            int
	primary             int
	secondary           int
	json1               int
	json2               int
	json3               int
	json4               int
	StyleParserColorMap map[string]termui.Color
}

var (
	theme *Theme
	once  sync.Once
)

func GetTheme() *Theme {

	once.Do(func() {
		theme = &Theme{1, 1, 1, 1, 1, 1, 1, 1, nil}
	})

	return theme
}

func (t *Theme) SetColors(base int, disabled int, primary int, secondary int, json1 int, json2 int, json3 int, json4 int) {
	t.base = base
	t.disabled = disabled
	t.primary = primary
	t.secondary = secondary
	t.json1 = json1
	t.json2 = json2
	t.json3 = json3
	t.json4 = json4

	t.StyleParserColorMap = map[string]termui.Color{}

	t.StyleParserColorMap["base"] = termui.Color(t.base)
	t.StyleParserColorMap["disabled"] = termui.Color(t.disabled)
	t.StyleParserColorMap["primary"] = termui.Color(t.primary)
	t.StyleParserColorMap["secondary"] = termui.Color(t.secondary)
	t.StyleParserColorMap["json1"] = termui.Color(t.json1)
	t.StyleParserColorMap["json2"] = termui.Color(t.json2)
	t.StyleParserColorMap["json3"] = termui.Color(t.json3)
	t.StyleParserColorMap["json4"] = termui.Color(t.json4)
}

func (t *Theme) GetColorByName(name string) termui.Color {
	return t.StyleParserColorMap[name]
}
