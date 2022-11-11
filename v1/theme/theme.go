package colors

import (
	"sync"

	termui "github.com/gizak/termui/v3"
)

type ColorScheme struct {
	id        string
	base      int
	disabled  int
	primary   int
	secondary int
	json1     int
	json2     int
	json3     int
	json4     int
	json5     int
}

//Theme uses xterm colors 1-255
type Theme struct {
	colorScheme         ColorScheme
	StyleParserColorMap map[string]termui.Color
}

var (
	theme   *Theme
	once    sync.Once
	Lavanda = ColorScheme{
		id:        "lavanda",
		base:      15,  //white
		disabled:  242, //gray
		primary:   91,
		secondary: 226,
		json1:     190,
		json2:     125,
		json3:     12,
		json4:     54,
		json5:     91,
	}
	MarromBombom = ColorScheme{
		id:        "marrombombom",
		base:      15,  //white
		disabled:  242, //gray
		primary:   236,
		secondary: 137,
		json1:     220,
		json2:     130,
		json3:     15,
		json4:     220,
		json5:     236,
	}
	Mono = ColorScheme{
		id:        "mono",
		base:      15, //white
		disabled:  242,
		primary:   234,
		secondary: 252,
		json1:     243,
		json2:     243,
		json3:     243,
		json4:     243,
		json5:     252,
		
	}
	Santa = ColorScheme{
		id:        "santa",
		base:      15, //white
		disabled:  242,
		primary:   196,
		secondary: 246,
		json1:     34,
		json2:     196,
		json3:     15,
		json4:     34,
		json5:     88,
		
	}

	themes map[string]ColorScheme
)

func GetManager() *Theme {

	once.Do(func() {
		themes = make(map[string]ColorScheme)
		themes[MarromBombom.id] = MarromBombom
		themes[Lavanda.id] = Lavanda
		themes[Mono.id] = Mono
		themes[Santa.id] = Santa
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
	t.StyleParserColorMap["json5"] = termui.Color(t.colorScheme.json5)
}

func (t *Theme) SetTheme(name string) {
	t.SetColors(themes[name])
}

func (t *Theme) GetColorByName(name string) termui.Color {
	return t.StyleParserColorMap[name]
}
