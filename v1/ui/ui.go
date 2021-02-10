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
	eventChannel = make(chan string)
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
	configflags  config.Config

	keysRegex, _    = regexp.Compile(`(\[37m)(.*?)(\[0m)`)
	stringsRegex, _ = regexp.Compile(`(\[32m)(.*?)(\[0m)`)
	numbersRegex, _ = regexp.Compile(`(\[36m)(.*?)(\[0m)`)
	booleanRegex, _ = regexp.Compile(`(\[33m)(.*?)(\[0m)`)
	termWidth       = 0
	termHeight      = 0
	formatter       = colorjson.NewFormatter()

	theme = t.GetManager()

	focused interface{}
)

func readEvents() {

	for {
		jsonBody := <-eventChannel

		eventList.Rows = append(eventList.Rows, fmt.Sprintf("%d %s", len(events), jsonBody))

		var obj map[string]interface{}
		json.Unmarshal([]byte(jsonBody), &obj)

		prettyJSON, _ := formatter.Marshal(obj)
		
		events = append(events, translateANSI(string(prettyJSON)))


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
		eventList.Title = fmt.Sprintf("%d/%d", _stats.Events, _stats.MaxEvents)
		stats.Text = fmt.Sprintf(" <%d>(fg:secondary) eps <%db>(fg:secondary) avg", _stats.Eps, _stats.AvgBytes)
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

	eventView.Rows = append(eventView.Rows, "")
	eventView.Rows = append(eventView.Rows, "")
	eventView.Rows = append(eventView.Rows, "")

	ui.Render(eventView)
}

func switchFocus() {

	if focused == eventList {
		focused = eventView
		eventList.BorderStyle.Fg = theme.GetColorByName("disabled")
		eventView.BorderStyle.Fg = theme.GetColorByName("base")
		eventList.TitleStyle.Fg = theme.GetColorByName("disabled")
		eventView.TitleStyle.Fg = theme.GetColorByName("base")
	} else {
		focused = eventList
		eventList.BorderStyle.Fg = theme.GetColorByName("base")
		eventView.BorderStyle.Fg = theme.GetColorByName("disabled")
		eventList.TitleStyle.Fg = theme.GetColorByName("base")
		eventView.TitleStyle.Fg = theme.GetColorByName("disabled")
	}

}

func fold(row int) {

	// still have to figure this out, styling messing up

	//eventView.Rows[row] = strings.Split(eventView.Rows[row], "\":")[1] + " {...}"
	//ui.Render(eventView)
}

func resize(width int, height int) {
	ui.Clear()

	grid.SetRect(0, 0, width, height-2)

	if configflags.Layout == "vertical" {

		grid.Set(
			ui.NewRow(0.3,
				ui.NewCol(1, eventList),
			),
			ui.NewRow(0.7,
				ui.NewCol(1, eventView),
			),
		)
	} else {
		grid.Set(
			ui.NewRow(1,
				ui.NewCol(0.2, eventList),
				ui.NewCol(0.8, eventView),
			),
		)
	}

	stats.SetRect(0, height-2, width-len(serverStatus.Text)-1, height-1)
	serverStatus.SetRect(width-len(serverStatus.Text)-1, height-2, width, height-1)

	footer.SetRect(0, height-1, width, height)
	version.SetRect(width-len(version.Text)-1, height-1, width, height)

	ui.Render(grid, stats, footer, version, serverStatus)
}

func Start(_config config.Config) {

	configflags = _config

	go server.Start(eventChannel, statsChannel, _config)

	if err := ui.Init(); err != nil {
		logs.Log(err)
	}
	defer ui.Close()

	formatter.Indent = 2

	t.GetManager().SetTheme("lavanda")
	eventView.Title = "JSON"
	eventView.WrapText = true

	footer.Text = " <q>(fg:disabled)uit <r>(fg:disabled)eset <l>(fg:disabled)ayout "
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

	serverStatus.Text = fmt.Sprintf("%s:%s", configflags.Host, configflags.Port)
	serverStatus.Border = false
	serverStatus.WrapText = false
	serverStatus.TextStyle.Fg = theme.GetColorByName("base")
	serverStatus.TextStyle.Bg = theme.GetColorByName("primary")

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
			case "<Space>":
				if focused == eventView {
					fold(eventView.SelectedRow)
				}
			case "r":
				reset()
				redrawFlag = true
			case "l":
				if configflags.Layout == "horizontal" {
					configflags.Layout = "vertical"
				} else {
					configflags.Layout = "horizontal"
				}
				resize(termWidth, termHeight)
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
