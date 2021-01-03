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
  "github.com/ugosan/logshark/v1/ansi8bit"
  "github.com/ugosan/logshark/v1/config"
  "github.com/ugosan/logshark/v1/logging"
  "github.com/ugosan/logshark/v1/server"
  logshark_widgets "github.com/ugosan/logshark/v1/widgets"
)

var (
  channel    = make(chan map[string]interface{})
  events     []string
  redrawFlag = true
  eventList  = widgets.NewList()
  eventView  = logshark_widgets.NewJSONView()
  footer     = logshark_widgets.NewFooter()
  stats      = logshark_widgets.NewFooter()
  version    = logshark_widgets.NewFooter()
  serverStatus = logshark_widgets.NewFooter()
  grid       = ui.NewGrid()
  logs       = logging.GetManager()

  keysRegex, _    = regexp.Compile(`(\[37m)(.*?)(\[0m)`)
  stringsRegex, _ = regexp.Compile(`(\[32m)(.*?)(\[0m)`)
  numbersRegex, _ = regexp.Compile(`(\[36m)(.*?)(\[0m)`)
  booleanRegex, _ = regexp.Compile(`(\[33m)(.*?)(\[0m)`)
  termWidth       = 0
  termHeight      = 0
  formatter       = colorjson.NewFormatter()

  focused interface{}
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

    ui.Render(eventView, eventList, footer, version, serverStatus)
    redrawFlag = false
  }
}

func reset() {
  events = events[:0]
  eventList.Rows = []string{}
  eventView.Rows = []string{}
  server.ResetStats()
}

func translateANSI(s string) string {

  s = stringsRegex.ReplaceAllString(s, "<$2>(fg:yellow2)")
  s = numbersRegex.ReplaceAllString(s, "<$2>(fg:darkviolet)")
  s = keysRegex.ReplaceAllString(s, "<$2>(fg:blue)")
  s = booleanRegex.ReplaceAllString(s, "<$2>(fg:purple4)")

  return s
}

func updateEventView() {

  eventView.Rows = strings.Split(events[eventList.SelectedRow], "\n")

  ui.Render(eventView)
}

func switchFocus() {

  if focused == eventList {
    focused = eventView
    eventList.BorderStyle.Fg = ansi8bit.Grey69
    eventView.BorderStyle.Fg = ansi8bit.White
    eventList.Title = "Events"
    eventView.Title = "JSON ●"
  } else {
    focused = eventList
    eventList.BorderStyle.Fg = ansi8bit.White
    eventView.BorderStyle.Fg = ansi8bit.Grey69
    eventList.Title = "Events ●"
    eventView.Title = "JSON"
  }

}

func updateStats() {

  stats.Text = fmt.Sprintf(" [%d](fg:white)/%d events %d e/s ", server.GetStats().Events, server.GetStats().MaxEvents, server.GetStats().Eps)
  ui.Render(stats)
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

  if err := ui.Init(); err != nil {
    logs.Log(err)
  }
  defer ui.Close()

  formatter.Indent = 2

  grid := ui.NewGrid()

  eventView.Title = "JSON"
  eventView.WrapText = true

  footer.Text = " [q](fg:yellow)uit [r](fg:yellow)eset"
  footer.Border = false
  footer.WrapText = false
  footer.TextStyle.Fg = ansi8bit.DarkViolet
  footer.TextStyle.Bg = ui.ColorWhite

  version.Text = "Logshark v1.0"
  version.Border = false
  version.WrapText = false
  version.TextStyle.Fg = ansi8bit.DarkViolet
  version.TextStyle.Bg = ui.ColorWhite

  stats.Border = false
  stats.WrapText = false
  stats.TextStyle.Fg = ui.ColorWhite
  stats.TextStyle.Bg = ansi8bit.DarkViolet

  serverStatus.Text = fmt.Sprintf("%s:%s", config.Host, config.Port)
  serverStatus.Border = false
  serverStatus.WrapText = false
  serverStatus.TextStyle.Fg = ui.ColorWhite
  serverStatus.TextStyle.Bg = ansi8bit.DarkViolet

  eventList.Title = "Events ●"
  eventList.TextStyle = ui.NewStyle(ansi8bit.Yellow1)
  eventList.WrapText = false

  eventList.BorderStyle.Fg = ansi8bit.White
  eventView.BorderStyle.Fg = ansi8bit.Grey69

  termWidth, termHeight = ui.TerminalDimensions()

  resize(termWidth, termHeight)

  focused = eventList

  go server.Start(channel, config)
  go readEvents()

  uiEvents := ui.PollEvents()
  ticker := time.NewTicker(time.Microsecond * 300).C

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
        ui.Render(grid)
      case "t":
        server.SendTestRequest()
        ui.Render(grid)
      case "<Resize>":
        payload := e.Payload.(ui.Resize)

        resize(payload.Width, payload.Height)
      }
    case <-ticker:
      updateStats()

      if len(eventList.Rows) == 1 {
        eventList.SelectedRow = 0
        updateEventView()
      }

      redraw()
    }
  }

}
