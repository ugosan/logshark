package main

import (
  "github.com/ugosan/logshark/cmd/server"
  "github.com/gdamore/tcell/v2"
  "github.com/rivo/tview"
  "time"
  "fmt"
  "encoding/json"
  "github.com/TylerBrock/colorjson"
)

var app = tview.NewApplication()
var eventList = tview.NewList()
var stats =  tview.NewTextView()
var eventView = tview.NewTextView()
var channel = make(chan map[string]interface{})


var formatter = colorjson.NewFormatter()

var events []interface{}

var redrawFlag = true

func read_events(){

  for {
    obj := <-channel
    
    events = append(events,obj)

    s, _ := json.Marshal(obj)
    
    eventList.AddItem(fmt.Sprintf("%d    %s", len(events), string(s)), "", 0, nil)
    redrawFlag = true
  }

}

//instead of redrawing at every event, redraws every 300 microseconds
func refresh() {
  for {

    if(redrawFlag){
      stats.SetText(fmt.Sprintf("%d/1000 | 0 e/s ", server.GetStats().Events))
      app.Draw()
      redrawFlag = false
    }

    time.Sleep(300 * time.Microsecond) 
  }
}

func navigation(event *tcell.EventKey) *tcell.EventKey {
  if event.Key() == tcell.KeyTab{

    if(eventList.HasFocus()){
      app.SetFocus(eventView)
    }else{
      app.SetFocus(eventList)
    }

    return nil
  }

  if event.Rune() == 't' {
    
    server.SendTestRequest()
  }

  if event.Rune() == 'r' {
    eventView.Clear()
    eventList.Clear()
    server.ResetStats()
  }

  if event.Rune() == 'q' {
    app.Stop()
  }
  return event
}

func Center(width, height int, p tview.Primitive) tview.Primitive {
  return tview.NewFlex().
    AddItem(nil, 0, 1, false).
    AddItem(tview.NewFlex().
      SetDirection(tview.FlexRow).
      AddItem(nil, 0, 1, false).
      AddItem(p, height, 1, true).
      AddItem(nil, 0, 1, false), width, 1, true).
    AddItem(nil, 0, 1, false)
}

func main() {

  go server.Start(channel)
  go read_events()
  go refresh()



  formatter.Indent = 4

  eventList.ShowSecondaryText(false)
  eventList.SetChangedFunc(func(line int, t string, t2 string, r rune) {
    eventView.Clear()

    s, _ := formatter.Marshal(events[line])
    fmt.Fprintf(eventView, "%s", tview.TranslateANSI(string(s)))
    
  })

  stats.
  SetDynamicColors(true).
  SetText(" 0/1000 | 0 e/s ")

  eventView.
    SetDynamicColors(true).
    SetRegions(true).
    SetWordWrap(true).
    SetBackgroundColor(tcell.ColorBlack)

  eventList.SetBorder(true)
  eventView.SetBorder(true)
  
  //app.SetInputCapture(navigation)
  eventView.SetInputCapture(navigation)
  eventList.SetInputCapture(navigation)



  title :=  tview.NewTextView().
  SetDynamicColors(true).
  SetText(" Logshark [gray]v0.1[white] ")
  

  footer :=  tview.NewTextView().
  SetDynamicColors(true).
  SetText(" [blue]r[white]efresh | [blue]s[white]ettings | [blue]q[white]uit")
  


  layout := tview.NewFlex().
  AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
    AddItem(title, 1, 1, false).
    AddItem(eventList, 0, 1, false).
    AddItem(Center(15,1,stats), 1, 1, false).
    AddItem(eventView, 0, 3, false).
    AddItem(footer, 1, 1, false), 0, 2, false)


  if err := app.SetRoot(layout, true).SetFocus(eventList).EnableMouse(true).Run(); err != nil {
    panic(err)
  }



}