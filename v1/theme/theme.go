package colors

import (
	"sync"

	termui "github.com/gizak/termui/v3"
)

type ColorScheme struct {
	base      int
	disabled  int
	primary   int
	secondary int
	json1     int
	json2     int
	json3     int
	json4     int
}

//Theme uses xterm colors 1-255
type Theme struct {
	colorScheme         ColorScheme
	StyleParserColorMap map[string]termui.Color
}

var (
	theme         *Theme
	once          sync.Once
	ThemeLavender = ColorScheme{base: 15, disabled: 242, primary: 91, secondary: 226, json1: 190, json2: 125, json3: 12, json4: 54}
)

func GetTheme() *Theme {

	once.Do(func() {
		colorScheme := ColorScheme{}
		theme = &Theme{colorScheme, nil}
	})

	return theme
}

func (t *Theme) SetColors(scheme ColorScheme) {
	t.colorScheme = scheme

	t.StyleParserColorMap = map[string]termui.Color{}

	t.StyleParserColorMap["base"] = termui.Color(t.colorScheme.base)
	t.StyleParserColorMap["disabled"] = termui.Color(t.colorScheme.disabled)
	t.StyleParserColorMap["primary"] = termui.Color(t.colorScheme.primary)
	t.StyleParserColorMap["secondary"] = termui.Color(t.colorScheme.secondary)
	t.StyleParserColorMap["json1"] = termui.Color(t.colorScheme.json1)
	t.StyleParserColorMap["json2"] = termui.Color(t.colorScheme.json2)
	t.StyleParserColorMap["json3"] = termui.Color(t.colorScheme.json3)
	t.StyleParserColorMap["json4"] = termui.Color(t.colorScheme.json4)
}

func (t *Theme) GetColorByName(name string) termui.Color {
	return t.StyleParserColorMap[name]
}