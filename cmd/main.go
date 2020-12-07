package main


import (
  "github.com/ugosan/logshark/cmd/server"
  "log"
  ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
  "fmt"
  "time"
  "encoding/json"
)
  
var channel = make(chan map[string]interface{})
var events []interface{}
var redrawFlag = true
var eventList = widgets.NewList()
var eventView = widgets.NewParagraph()
var footer = widgets.NewParagraph()
var grid = ui.NewGrid()

func read_events(){

  for {
    obj := <-channel
    events = append(events, obj)
    s, _ := json.Marshal(obj)
    eventList.Rows = append(eventList.Rows, fmt.Sprintf("%d %s", len(events), string(s)))

    footer.Text = fmt.Sprintf("%d/1000 | 0 e/s ", server.GetStats().Events)
    ui.Render(grid)
  }

}

//instead of redrawing at every event, redraws every 300 microseconds
func redraw() {
  for {

    if(redrawFlag){
      footer.Text = fmt.Sprintf("%d/1000 | 0 e/s ", server.GetStats().Events)
      ui.Render(grid)
      redrawFlag = false
    }

    time.Sleep(300 * time.Microsecond) 
  }
}

func reset() {
  events = events[:0]
  eventList.Rows = []string{}
  eventView.Text = ""
  server.ResetStats()
}

func updateEventView() {

  if(eventList.SelectedRow>-1){
    s, _ := json.MarshalIndent(events[eventList.SelectedRow], "", "  ")

    //stringsRegex, _ := regexp.Compile(`(".*")(: )`)

    //pretty := stringsRegex.ReplaceAllString(string(s), "[$1](fg:yellow): ")
  
    eventView.Text = string(s)

    ui.Render(eventList, eventView)
  }

}

func main() {

  go server.Start(channel)

  if err := ui.Init(); err != nil {
    log.Fatalf("failed to initialize termui: %v", err)
  }
  defer ui.Close()


  eventView.Title = "preview"
  footer.Border = true
  
  tsl := widgets.NewSparkline()
  tsl.LineColor = ui.ColorGreen
  tsl.Data = []float64{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
  tsls := widgets.NewSparklineGroup(tsl)
  //tsls.Title = "  "
  tsls.BorderStyle.Fg = ui.ColorWhite

  eventList.Title = "List"
  eventList.TextStyle = ui.NewStyle(ui.ColorYellow)
  eventList.WrapText = false

  grid := ui.NewGrid()
  termWidth, termHeight := ui.TerminalDimensions()
  grid.SetRect(0, 0, termWidth, termHeight)

  grid.Set(
    ui.NewRow(0.92,
      ui.NewCol(0.2, eventList),
      ui.NewCol(0.8, eventView),
    ),

    ui.NewRow(0.08,
      ui.NewCol(0.2, tsls),
      ui.NewCol(0.8, footer),
    ),

  )

  ui.Render(grid)
  
  
  go read_events()
  go redraw()

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
        grid.SetRect(0, 0, payload.Width, payload.Height)
        ui.Clear()
        ui.Render(grid)
      }
    case <-ticker:
      ui.Render(eventList, eventView, footer)
    }
  }

}