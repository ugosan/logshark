package ui


import (
	"github.com/ugosan/logshark/cmd/server"
  "github.com/ugosan/logshark/cmd/config"
  logshark_widgets "github.com/ugosan/logshark/cmd/widgets"
  "github.com/TylerBrock/colorjson"
  "log"
  ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
  "fmt"
  "time"
  "regexp"
  "encoding/json"
)
  
var (
  channel = make(chan map[string]interface{})
  events []interface{}
  redrawFlag = true
  eventList = widgets.NewList()
  eventView = logshark_widgets.NewParagraph()
  footer = widgets.NewParagraph()
  grid = ui.NewGrid()
)


func read_events(){

  for {
    obj := <-channel
    events = append(events, obj)
    s, _ := json.Marshal(obj)
    eventList.Rows = append(eventList.Rows, fmt.Sprintf("%d %s", len(events), string(s)))

    footer.Text = fmt.Sprintf("%d/1000 | %d e/s ", server.GetStats().Events, server.GetStats().Eps)
    ui.Render(grid)
  }

}

//instead of redrawing at every event, redraws every 300 microseconds
func redraw() {
  for {

    if(redrawFlag){
      
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

    f := colorjson.NewFormatter()
    f.Indent = 2

    s, _ := f.Marshal(events[eventList.SelectedRow])
    //fmt.Println(string(s))

    pretty := string(s)
    keysRegex, _ := regexp.Compile(`(\[37m)(.*?)(\[0m)`)
    stringsRegex, _ := regexp.Compile(`(\[32m)(.*?)(\[0m)`)
    numbersRegex, _ := regexp.Compile(`(\[36m)(.*?)(\[0m)`)
    booleanRegex, _ := regexp.Compile(`(\[33m)(.*?)(\[0m)`)


    pretty = stringsRegex.ReplaceAllString(pretty, "<$2>(fg:yellow)")
    pretty = numbersRegex.ReplaceAllString(pretty, "<$2>(fg:magenta)")
    pretty = keysRegex.ReplaceAllString(pretty, "<$2>(fg:blue)")
    pretty = booleanRegex.ReplaceAllString(pretty, "<$2>(fg:green)")

    //s, _ := json.MarshalIndent(events[eventList.SelectedRow], "", "  ")

    eventView.Text = pretty

    ui.Render(eventList, eventView)
  }

}

func updateStats() {
  footer.Text = fmt.Sprintf("%d/1000 | %d e/s ", server.GetStats().Events, server.GetStats().Eps)
}


func Start(config config.Config) {

  go server.Start(channel, config)

  if err := ui.Init(); err != nil {
    log.Fatalf("failed to initialize termui: %v", err)
  }
  defer ui.Close()


  eventView.Title = "preview"
  footer.Border = true
  
  tsl := widgets.NewSparkline()
  tsl.LineColor = ui.ColorGreen
  tsl.Data = []float64{2, 2, 1, 2, 3, 2, 1, 2, 2, 2, 2, 2, 2, 3, 1, 4, 5, 2, 2, 5}
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
      case "u":
        eventView.Text = "[3] [\"color\"](fg:white,bg:green) output"
        
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
      updateStats()
      if(len(eventList.Rows) == 1){
        eventList.SelectedRow = 0
        updateEventView()
      }
    }
  }

}