package main


import (
  "github.com/ugosan/logshark/cmd/server"
  "log"
  ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
  "fmt"
  "encoding/json"
  "github.com/TylerBrock/colorjson"
)

var channel = make(chan map[string]interface{})
var events []interface{}

var formatter = colorjson.NewFormatter()
var eventList = widgets.NewList()
var eventView = widgets.NewParagraph()


func read_events(){

  for {
    obj := <-channel
    events = append(events, obj)
    s, _ := json.Marshal(obj)
    eventList.Rows = append(eventList.Rows, fmt.Sprintf("%d %s", len(events), string(s)))
  }

}

func reset() {
  events = events[:0]
  eventList.Rows = []string{}
  eventView.Text = ""
  server.ResetStats()
}

func updateEventView() {

  formatter.Indent = 2
  s, _ := formatter.Marshal(events[eventList.SelectedRow])


  eventView.Text = string(s)

}

func main() {

  go server.Start(channel)
  go read_events()



  if err := ui.Init(); err != nil {
    log.Fatalf("failed to initialize termui: %v", err)
  }
  defer ui.Close()


  eventView.Title = "preview"

  footer := widgets.NewParagraph()
  footer.Text = "lalalalala lalalala"
  footer.Border = false


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
      ui.NewCol(1, footer),
    ),

  )

  ui.Render(grid)


  for {
    select {
      case e := <-ui.PollEvents():
        switch e.ID {
        case "q", "<C-c>":
          return
        case "<Down>":
          eventList.ScrollDown()
          updateEventView()
          ui.Render(grid)
        case "<Up>":
          eventList.ScrollUp()
          updateEventView()
          ui.Render(grid)
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
    }
  }
}