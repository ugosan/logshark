package ui

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/ugosan/logshark/v1/config"
	"github.com/ugosan/logshark/v1/logging"
	"github.com/ugosan/logshark/v1/server"
	t "github.com/ugosan/logshark/v1/theme"
	logshark_widgets "github.com/ugosan/logshark/v1/widgets"
)

var (
	eventChannel = make(chan map[string]interface{})
	statsChannel = make(chan server.Stats)
	events       []string
	redrawFlag   = true
	eventList    = widgets.NewList()
	eventView    = logshark_widgets.NewJSONView()
	footer       = logshark_widgets.NewFooter()
	stats        = logshark_widgets.NewFooter()
	version      = logshark_widgets.NewFooter()
	serverStatus = logshark_widgets.NewFooter()
	grid         = ui.NewGrid()
	logs         = logging.GetManager()

	keysRegex, _    = regexp.Compile(`(\[37m)(.*?)(\[0m)`)
	stringsRegex, _ = regexp.Compile(`(\[32m)(.*?)(\[0m)`)
	numbersRegex, _ = regexp.Compile(`(\[36m)(.*?)(\[0m)`)
	booleanRegex, _ = regexp.Compile(`(\[33m)(.*?)(\[0m)`)
	termWidth       = 0
	termHeight      = 0
	formatter       = colorjson.NewFormatter()

	theme = t.GetTheme()

	focused interface{}
)

func readEvents() {

	for {
		obj := <-eventChannel

		prettyJSON, _ := formatter.Marshal(obj)
		events = append(events, translateANSI(string(prettyJSON)))

		s, _ := json.Marshal(obj)
		eventList.Rows = append(eventList.Rows, fmt.Sprintf("%d %s", len(events), string(s)))

		redrawFlag = true

		if len(eventList.Rows) == 1 {
			eventList.SelectedRow = 0
			updateEventView()
		}

	}

}

//instead of redrawing at every event, redraws every 300 microseconds
func redraw() {
	if redrawFlag == true {

		ui.Render(eventView, eventList, footer, stats, version, serverStatus)
		redrawFlag = false
	}
}

func readStats() {

	for {
		_stats := <-statsChannel
		stats.Text = fmt.Sprintf(" <%d>(fg:base)/%d events %d e/s ", _stats.Events, _stats.MaxEvents, _stats.Eps)
		redrawFlag = true
	}

}

func reset() {
	events = events[:0]
	eventList.Rows = []string{}
	eventView.Rows = []string{}
	server.ResetStats()
}

func translateANSI(s string) string {

	s = stringsRegex.ReplaceAllString(s, "<$2>(fg:json1)")
	s = numbersRegex.ReplaceAllString(s, "<$2>(fg:json2)")
	s = keysRegex.ReplaceAllString(s, "<$2>(fg:json3)")
	s = booleanRegex.ReplaceAllString(s, "<$2>(fg:json4)")

	return s
}

func updateEventView() {

	eventView.Rows = strings.Split(events[eventList.SelectedRow], "\n")

	ui.Render(eventView)
}

func switchFocus() {

	if focused == eventList {
		focused = eventView
		eventList.BorderStyle.Fg = theme.GetColorByName("disabled")
		eventView.BorderStyle.Fg = theme.GetColorByName("base")
		eventList.Title = "Events"
		eventView.Title = "JSON ●"
	} else {
		focused = eventList
		eventList.BorderStyle.Fg = theme.GetColorByName("base")
		eventView.BorderStyle.Fg = theme.GetColorByName("disabled")
		eventList.Title = "Events ●"
		eventView.Title = "JSON"
	}

}

func resize(width int, height int) {
	ui.Clear()

	grid.SetRect(0, 0, width, height-2)

	grid.Set(
		ui.NewRow(1,
			ui.NewCol(0.2, eventList),
			ui.NewCol(0.8, eventView),
		),
	)

	stats.SetRect(0, height-2, width-len(serverStatus.Text)-1, height-1)
	serverStatus.SetRect(width-len(serverStatus.Text)-1, height-2, width, height-1)

	footer.SetRect(0, height-1, width, height)
	version.SetRect(width-len(version.Text)-1, height-1, width, height)

	ui.Render(grid, stats, footer, version, serverStatus)
}

func Start(config config.Config) {

	go server.Start(eventChannel, statsChannel, config)

	if err := ui.Init(); err != nil {
		logs.Log(err)
	}
	defer ui.Close()

	formatter.Indent = 2

	theme.SetColors(15, 242, 91, 226, 190, 125, 12, 54)

	eventView.Title = "JSON"
	eventView.WrapText = true

	footer.Text = " <q>(fg:primary)uit <r>(fg:primary)eset"
	footer.Border = false
	footer.WrapText = false
	footer.TextStyle.Fg = theme.GetColorByName("primary")
	footer.TextStyle.Bg = theme.GetColorByName("base")

	version.Text = "Logshark v1.0"
	version.Border = false
	version.WrapText = false
	version.TextStyle.Fg = theme.GetColorByName("primary")
	version.TextStyle.Bg = theme.GetColorByName("base")

	stats.Border = false
	stats.WrapText = false
	stats.TextStyle.Fg = theme.GetColorByName("base")
	stats.TextStyle.Bg = theme.GetColorByName("primary")

	serverStatus.Text = fmt.Sprintf("%s:%s", config.Host, config.Port)
	serverStatus.Border = false
	serverStatus.WrapText = false
	serverStatus.TextStyle.Fg = theme.GetColorByName("base")
	serverStatus.TextStyle.Bg = theme.GetColorByName("primary")

	eventList.Title = "Events ●"
	eventList.TextStyle.Fg = theme.GetColorByName("secondary")
	eventList.WrapText = false

	eventList.BorderStyle.Fg = theme.GetColorByName("base")
	eventView.BorderStyle.Fg = theme.GetColorByName("disabled")

	termWidth, termHeight = ui.TerminalDimensions()

	resize(termWidth, termHeight)

	focused = eventList

	redrawFlag = true

	go readEvents()
	go readStats()
	go redraw()

	uiEvents := ui.PollEvents()

	redrawTicker := time.NewTicker(time.Microsecond * 300).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Down>":
				if focused == eventList {
					if len(events) > 0 {
						eventList.ScrollDown()
						updateEventView()
						ui.Render(eventList)
					}
				} else {
					if len(eventView.Rows) > 0 {
						eventView.ScrollDown()
						ui.Render(eventView)
					}
				}
			case "<Up>":
				if focused == eventList {
					if len(events) > 0 {
						eventList.ScrollUp()
						updateEventView()
						ui.Render(eventList)
					}
				} else {
					if len(eventView.Rows) > 0 {
						eventView.ScrollUp()
						ui.Render(eventView)
					}
				}
			case "<Tab>":
				switchFocus()
				ui.Render(eventList, eventView)
			case "r":
				reset()
				redrawFlag = true
			case "t":
				server.SendTestRequest()
			case "<Resize>":
				payload := e.Payload.(ui.Resize)

				resize(payload.Width, payload.Height)
			}
		case <-redrawTicker:
			redraw()
		}

	}

}
