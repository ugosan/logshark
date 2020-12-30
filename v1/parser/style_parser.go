package parser

import (
	"strings"
	. "github.com/gizak/termui/v3"
	"github.com/ugosan/logshark/v1/ansi8bit"
	
)

const (
	tokenFg       = "fg"
	tokenBg       = "bg"
	tokenModifier = "mod"

	tokenItemSeparator  = ","
	tokenValueSeparator = ":"

	tokenBeginStyledText = '<'
	tokenEndStyledText   = '>'

	tokenBeginStyle = '('
	tokenEndStyle   = ')'
)

type parserState uint

const (
	parserStateDefault parserState = iota
	parserStateStyleItems
	parserStateStyledText
	ColorDarkRed = 88
)

// StyleParserColorMap can be modified to add custom color parsing to text
var StyleParserColorMap = map[string]Color{
	"black": ansi8bit.Black,
	"maroon": ansi8bit.Maroon,
	"green": ansi8bit.Green,
	"olive": ansi8bit.Olive,
	"navy": ansi8bit.Navy,
	"purple": ansi8bit.Purple,
	"teal": ansi8bit.Teal,
	"silver": ansi8bit.Silver,
	"grey": ansi8bit.Grey,
	"red": ansi8bit.Red,
	"lime": ansi8bit.Lime,
	"yellow": ansi8bit.Yellow,
	"blue": ansi8bit.Blue,
	"fuchsia": ansi8bit.Fuchsia,
	"aqua": ansi8bit.Aqua,
	"white": ansi8bit.White,
	"grey0": ansi8bit.Grey0,
	"navyblue": ansi8bit.NavyBlue,
	"darkblue": ansi8bit.DarkBlue,
	"blue3": ansi8bit.Blue3,
	"blue1": ansi8bit.Blue1,
	"darkgreen": ansi8bit.DarkGreen,
	"deepskyblue4": ansi8bit.DeepSkyBlue4,
	"dodgerblue3": ansi8bit.DodgerBlue3,
	"dodgerblue2": ansi8bit.DodgerBlue2,
	"green4": ansi8bit.Green4,
	"springgreen4": ansi8bit.SpringGreen4,
	"turquoise4": ansi8bit.Turquoise4,
	"deepskyblue3": ansi8bit.DeepSkyBlue3,
	"dodgerblue1": ansi8bit.DodgerBlue1,
	"green3": ansi8bit.Green3,
	"springgreen3": ansi8bit.SpringGreen3,
	"darkcyan": ansi8bit.DarkCyan,
	"lightseagreen": ansi8bit.LightSeaGreen,
	"deepskyblue2": ansi8bit.DeepSkyBlue2,
	"deepskyblue1": ansi8bit.DeepSkyBlue1,
	"springgreen2": ansi8bit.SpringGreen2,
	"cyan3": ansi8bit.Cyan3,
	"darkturquoise": ansi8bit.DarkTurquoise,
	"turquoise2": ansi8bit.Turquoise2,
	"green1": ansi8bit.Green1,
	"springgreen1": ansi8bit.SpringGreen1,
	"mediumspringgreen": ansi8bit.MediumSpringGreen,
	"cyan2": ansi8bit.Cyan2,
	"cyan1": ansi8bit.Cyan1,
	"darkred": ansi8bit.DarkRed,
	"deeppink4": ansi8bit.DeepPink4,
	"purple4": ansi8bit.Purple4,
	"purple3": ansi8bit.Purple3,
	"blueviolet": ansi8bit.BlueViolet,
	"orange4": ansi8bit.Orange4,
	"grey37": ansi8bit.Grey37,
	"mediumpurple4": ansi8bit.MediumPurple4,
	"slateblue3": ansi8bit.SlateBlue3,
	"royalblue1": ansi8bit.RoyalBlue1,
	"chartreuse4": ansi8bit.Chartreuse4,
	"darkseagreen4": ansi8bit.DarkSeaGreen4,
	"paleturquoise4": ansi8bit.PaleTurquoise4,
	"steelblue": ansi8bit.SteelBlue,
	"steelblue3": ansi8bit.SteelBlue3,
	"cornflowerblue": ansi8bit.CornflowerBlue,
	"chartreuse3": ansi8bit.Chartreuse3,
	"cadetblue": ansi8bit.CadetBlue,
	"skyblue3": ansi8bit.SkyBlue3,
	"steelblue1": ansi8bit.SteelBlue1,
	"palegreen3": ansi8bit.PaleGreen3,
	"seagreen3": ansi8bit.SeaGreen3,
	"aquamarine3": ansi8bit.Aquamarine3,
	"mediumturquoise": ansi8bit.MediumTurquoise,
	"chartreuse2": ansi8bit.Chartreuse2,
	"seagreen2": ansi8bit.SeaGreen2,
	"seagreen1": ansi8bit.SeaGreen1,
	"aquamarine1": ansi8bit.Aquamarine1,
	"darkslategray2": ansi8bit.DarkSlateGray2,
	"darkmagenta": ansi8bit.DarkMagenta,
	"darkviolet": ansi8bit.DarkViolet,
	"lightpink4": ansi8bit.LightPink4,
	"plum4": ansi8bit.Plum4,
	"mediumpurple3": ansi8bit.MediumPurple3,
	"slateblue1": ansi8bit.SlateBlue1,
	"yellow4": ansi8bit.Yellow4,
	"wheat4": ansi8bit.Wheat4,
	"grey53": ansi8bit.Grey53,
	"lightslategrey": ansi8bit.LightSlateGrey,
	"mediumpurple": ansi8bit.MediumPurple,
	"lightslateblue": ansi8bit.LightSlateBlue,
	"darkolivegreen3": ansi8bit.DarkOliveGreen3,
	"darkseagreen": ansi8bit.DarkSeaGreen,
	"lightskyblue3": ansi8bit.LightSkyBlue3,
	"skyblue2": ansi8bit.SkyBlue2,
	"darkseagreen3": ansi8bit.DarkSeaGreen3,
	"darkslategray3": ansi8bit.DarkSlateGray3,
	"skyblue1": ansi8bit.SkyBlue1,
	"chartreuse1": ansi8bit.Chartreuse1,
	"lightgreen": ansi8bit.LightGreen,
	"palegreen1": ansi8bit.PaleGreen1,
	"darkslategray1": ansi8bit.DarkSlateGray1,
	"red3": ansi8bit.Red3,
	"mediumvioletred": ansi8bit.MediumVioletRed,
	"magenta3": ansi8bit.Magenta3,
	"darkorange3": ansi8bit.DarkOrange3,
	"indianred": ansi8bit.IndianRed,
	"hotpink3": ansi8bit.HotPink3,
	"mediumorchid3": ansi8bit.MediumOrchid3,
	"mediumorchid": ansi8bit.MediumOrchid,
	"mediumpurple2": ansi8bit.MediumPurple2,
	"darkgoldenrod": ansi8bit.DarkGoldenrod,
	"lightsalmon3": ansi8bit.LightSalmon3,
	"rosybrown": ansi8bit.RosyBrown,
	"grey63": ansi8bit.Grey63,
	"mediumpurple1": ansi8bit.MediumPurple1,
	"gold3": ansi8bit.Gold3,
	"darkkhaki": ansi8bit.DarkKhaki,
	"navajowhite3": ansi8bit.NavajoWhite3,
	"grey69": ansi8bit.Grey69,
	"lightsteelblue3": ansi8bit.LightSteelBlue3,
	"lightsteelblue": ansi8bit.LightSteelBlue,
	"yellow3": ansi8bit.Yellow3,
	"darkseagreen2": ansi8bit.DarkSeaGreen2,
	"lightcyan3": ansi8bit.LightCyan3,
	"lightskyblue1": ansi8bit.LightSkyBlue1,
	"greenyellow": ansi8bit.GreenYellow,
	"darkolivegreen2": ansi8bit.DarkOliveGreen2,
	"darkseagreen1": ansi8bit.DarkSeaGreen1,
	"paleturquoise1": ansi8bit.PaleTurquoise1,
	"deeppink3": ansi8bit.DeepPink3,
	"magenta2": ansi8bit.Magenta2,
	"hotpink2": ansi8bit.HotPink2,
	"orchid": ansi8bit.Orchid,
	"mediumorchid1": ansi8bit.MediumOrchid1,
	"orange3": ansi8bit.Orange3,
	"lightpink3": ansi8bit.LightPink3,
	"pink3": ansi8bit.Pink3,
	"plum3": ansi8bit.Plum3,
	"violet": ansi8bit.Violet,
	"lightgoldenrod3": ansi8bit.LightGoldenrod3,
	"tan": ansi8bit.Tan,
	"mistyrose3": ansi8bit.MistyRose3,
	"thistle3": ansi8bit.Thistle3,
	"plum2": ansi8bit.Plum2,
	"khaki3": ansi8bit.Khaki3,
	"lightgoldenrod2": ansi8bit.LightGoldenrod2,
	"lightyellow3": ansi8bit.LightYellow3,
	"grey84": ansi8bit.Grey84,
	"lightsteelblue1": ansi8bit.LightSteelBlue1,
	"yellow2": ansi8bit.Yellow2,
	"darkolivegreen1": ansi8bit.DarkOliveGreen1,
	"honeydew2": ansi8bit.Honeydew2,
	"lightcyan1": ansi8bit.LightCyan1,
	"red1": ansi8bit.Red1,
	"deeppink2": ansi8bit.DeepPink2,
	"deeppink1": ansi8bit.DeepPink1,
	"magenta1": ansi8bit.Magenta1,
	"orangered1": ansi8bit.OrangeRed1,
	"indianred1": ansi8bit.IndianRed1,
	"hotpink": ansi8bit.HotPink,
	"darkorange": ansi8bit.DarkOrange,
	"salmon1": ansi8bit.Salmon1,
	"lightcoral": ansi8bit.LightCoral,
	"palevioletred1": ansi8bit.PaleVioletRed1,
	"orchid2": ansi8bit.Orchid2,
	"orchid1": ansi8bit.Orchid1,
	"orange1": ansi8bit.Orange1,
	"sandybrown": ansi8bit.SandyBrown,
	"lightsalmon1": ansi8bit.LightSalmon1,
	"lightpink1": ansi8bit.LightPink1,
	"pink1": ansi8bit.Pink1,
	"plum1": ansi8bit.Plum1,
	"gold1": ansi8bit.Gold1,
	"navajowhite1": ansi8bit.NavajoWhite1,
	"mistyrose1": ansi8bit.MistyRose1,
	"thistle1": ansi8bit.Thistle1,
	"yellow1": ansi8bit.Yellow1,
	"lightgoldenrod1": ansi8bit.LightGoldenrod1,
	"khaki1": ansi8bit.Khaki1,
	"wheat1": ansi8bit.Wheat1,
	"cornsilk1": ansi8bit.Cornsilk1,
	"grey100": ansi8bit.Grey100,
	"grey3": ansi8bit.Grey3,
	"grey7": ansi8bit.Grey7,
	"grey11": ansi8bit.Grey11,
	"grey15": ansi8bit.Grey15,
	"grey19": ansi8bit.Grey19,
	"grey23": ansi8bit.Grey23,
	"grey27": ansi8bit.Grey27,
	"grey30": ansi8bit.Grey30,
	"grey35": ansi8bit.Grey35,
	"grey39": ansi8bit.Grey39,
	"grey42": ansi8bit.Grey42,
	"grey46": ansi8bit.Grey46,
	"grey50": ansi8bit.Grey50,
	"grey54": ansi8bit.Grey54,
	"grey58": ansi8bit.Grey58,
	"grey62": ansi8bit.Grey62,
	"grey66": ansi8bit.Grey66,
	"grey70": ansi8bit.Grey70,
	"grey74": ansi8bit.Grey74,
	"grey78": ansi8bit.Grey78,
	"grey82": ansi8bit.Grey82,
	"grey85": ansi8bit.Grey85,
	"grey89": ansi8bit.Grey89,
	"grey93": ansi8bit.Grey93,
}

var modifierMap = map[string]Modifier{
	"bold":      ModifierBold,
	"underline": ModifierUnderline,
	"reverse":   ModifierReverse,
}

// readStyle translates an []rune like `fg:red,mod:bold,bg:white` to a style
func readStyle(runes []rune, defaultStyle Style) Style {
	style := defaultStyle
	split := strings.Split(string(runes), tokenItemSeparator)
	for _, item := range split {
		pair := strings.Split(item, tokenValueSeparator)
		if len(pair) == 2 {
			switch pair[0] {
			case tokenFg:
				style.Fg = StyleParserColorMap[pair[1]]
			case tokenBg:
				style.Bg = StyleParserColorMap[pair[1]]
			case tokenModifier:
				style.Modifier = modifierMap[pair[1]]
			}
		}
	}
	return style
}

// ParseStyles parses a string for embedded Styles and returns []Cell with the correct styling.
// Uses defaultStyle for any text without an embedded style.
// Syntax is of the form [text](fg:<color>,mod:<attribute>,bg:<color>).
// Ordering does not matter. All fields are optional.
func ParseStyles(s string, defaultStyle Style) []Cell {
	cells := []Cell{}
	runes := []rune(s)
	state := parserStateDefault
	styledText := []rune{}
	styleItems := []rune{}
	squareCount := 0

	reset := func() {
		styledText = []rune{}
		styleItems = []rune{}
		state = parserStateDefault
		squareCount = 0
	}

	rollback := func() {
		cells = append(cells, RunesToStyledCells(styledText, defaultStyle)...)
		cells = append(cells, RunesToStyledCells(styleItems, defaultStyle)...)
		reset()
	}

	// chop first and last runes
	chop := func(s []rune) []rune {
		return s[1 : len(s)-1]
	}

	for i, _rune := range runes {
		switch state {
		case parserStateDefault:
			if _rune == tokenBeginStyledText {
				state = parserStateStyledText
				squareCount = 1
				styledText = append(styledText, _rune)
			} else {
				cells = append(cells, Cell{_rune, defaultStyle})
			}
		case parserStateStyledText:
			switch {
			case squareCount == 0:
				switch _rune {
				case tokenBeginStyle:
					state = parserStateStyleItems
					styleItems = append(styleItems, _rune)
				default:
					rollback()
					switch _rune {
					case tokenBeginStyledText:
						state = parserStateStyledText
						squareCount = 1
						styleItems = append(styleItems, _rune)
					default:
						cells = append(cells, Cell{_rune, defaultStyle})
					}
				}
			case len(runes) == i+1:
				rollback()
				styledText = append(styledText, _rune)
			case _rune == tokenBeginStyledText:
				squareCount++
				styledText = append(styledText, _rune)
			case _rune == tokenEndStyledText:
				squareCount--
				styledText = append(styledText, _rune)
			default:
				styledText = append(styledText, _rune)
			}
		case parserStateStyleItems:
			styleItems = append(styleItems, _rune)
			if _rune == tokenEndStyle {
				style := readStyle(chop(styleItems), defaultStyle)
				cells = append(cells, RunesToStyledCells(chop(styledText), style)...)
				reset()
			} else if len(runes) == i+1 {
				rollback()
			}
		}
	}

	return cells
}