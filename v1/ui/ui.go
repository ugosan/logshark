package ui

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/TylerBrock/colorjson"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/ugosan/logshark/v1/config"
	"github.com/ugosan/logshark/v1/logging"
	"github.com/ugosan/logshark/v1/server"
	logshark_widgets "github.com/ugosan/logshark/v1/widgets"
)

var (
  channel = make(chan map[string]interface{})
  events []string
  redrawFlag = true
  eventList = widgets.NewList()
  eventView = logshark_widgets.NewParagraph()
  footer = logshark_widgets.NewFooter()
  stats = logshark_widgets.NewFooter()
  grid = ui.NewGrid()
  logs = logging.GetManager()
  
  

  keysRegex, _ = regexp.Compile(`(\[37m)(.*?)(\[0m)`)
  stringsRegex, _ = regexp.Compile(`(\[32m)(.*?)(\[0m)`)
  numbersRegex, _ = regexp.Compile(`(\[36m)(.*?)(\[0m)`)
  booleanRegex, _ = regexp.Compile(`(\[33m)(.*?)(\[0m)`)
  termWidth = 0
  termHeight = 0
  formatter = colorjson.NewFormatter()
)

func readEvents() {

	for {
		obj := <-channel

		prettyJson, _ := formatter.Marshal(obj)
		events = append(events, translateANSI(string(prettyJson)))

		s, _ := json.Marshal(obj)
		eventList.Rows = append(eventList.Rows, fmt.Sprintf("%d %s", len(events), string(s)))

		redrawFlag = true
	}

}

//instead of redrawing at every event, redraws every 300 microseconds
func redraw() {
	if redrawFlag {

		ui.Render(eventView, eventList, footer)
		redrawFlag = false
	}
}

func reset() {
	events = events[:0]
	eventList.Rows = []string{}
	eventView.Text = ""
	server.ResetStats()
}

func translateANSI(s string) string {

	s = stringsRegex.ReplaceAllString(s, "<$2>(fg:yellow)")
	s = numbersRegex.ReplaceAllString(s, "<$2>(fg:magenta)")
	s = keysRegex.ReplaceAllString(s, "<$2>(fg:blue)")
	s = booleanRegex.ReplaceAllString(s, "<$2>(fg:green)")

	return s
}

func updateEventView() {

	if eventList.SelectedRow > -1 {

		eventView.Text = events[eventList.SelectedRow]

		ui.Render(eventList, eventView)
	}

}

func updateStats() {

  statsText := fmt.Sprintf(" %d/%d | %d e/s ", server.GetStats().Events, server.GetStats().MaxEvents, server.GetStats().Eps)

  stats.Text = statsText
  ui.Render(stats)
}

func Start(config config.Config) {

  go server.Start(channel, config)

  if err := ui.Init(); err != nil {
    logs.Log(err)
  }
  defer ui.Close()

  formatter.Indent = 2
  eventView.Title = "preview"
  footer.Border = false
  footer.WrapText = false
  footer.TextStyle.Fg = ui.ColorBlack
  footer.TextStyle.Bg = ui.ColorWhite


  stats.Border = false
  stats.WrapText = false
  stats.TextStyle.Fg = ui.ColorWhite
  stats.TextStyle.Bg = ui.ColorBlack

  eventList.Title = "List"
  eventList.TextStyle = ui.NewStyle(ui.ColorYellow)
  eventList.WrapText = false

  grid := ui.NewGrid()
  termWidth, termHeight = ui.TerminalDimensions()
  
  grid.SetRect(0, 0, termWidth, termHeight-2)

  grid.Set(
    ui.NewRow(1,
      ui.NewCol(0.2, eventList),
      ui.NewCol(0.8, eventView),
    ),
  )

  footer.SetRect(0, termHeight-1, termWidth, termHeight)
  stats.SetRect(0, termHeight-2, termWidth, termHeight-1)

  footer.Text = " [q](fg:yellow)uit [r](fg:yellow)eset"

  ui.Render(grid)
  
  go readEvents()

  uiEvents := ui.PollEvents()
  ticker := time.NewTicker(time.Microsecond*300).C

  for {
    select {
    case e := <-uiEvents:
      switch e.ID {
      case "q", "<C-c>":
        return
      case "<Down>":
        eventList.ScrollDown()
        updateEventView()
      case "<Up>":
        eventList.ScrollUp()
        updateEventView()
      case "r":
        reset()
        ui.Render(grid)
      case "t":
        server.SendTestRequest()
        ui.Render(grid)
      case "<Resize>":
        payload := e.Payload.(ui.Resize)

        termWidth = payload.Width
        termHeight = payload.Height
        
        grid.SetRect(0, 0, termWidth, termHeight-1)
        footer.SetRect(0, payload.Height-1, payload.Width, payload.Height)
        ui.Clear()
        ui.Render(grid)
      }
    case <-ticker:
      updateStats()
      if(len(eventList.Rows) == 1){
        eventList.SelectedRow = 0
        updateEventView()
      }
      
      redraw()
    }
  }

}
