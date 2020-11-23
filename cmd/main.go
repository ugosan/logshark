// Demo code for the List primitive.
package main

import (
  "github.com/ugosan/logshark/cmd/server"
  "github.com/rivo/tview"
  "time"
  "fmt"
  "encoding/json"
  "github.com/TylerBrock/colorjson"
)

var eventList = tview.NewList()
var c = make(chan map[string]interface{})
var app = tview.NewApplication()
var eventView = tview.NewTextView()

func read_events(){

  f := colorjson.NewFormatter()
  f.Indent = 4
    
  for {
    obj := <-c
    
    //s, _ := f.Marshal(obj)
    //fmt.Println(string(s))
    //j, err := json.MarshalIndent(obj, "", "")

    s, _ := json.Marshal(obj)
    
    eventList.AddItem(string(s), "", 0, nil)
    
    
    //if err != nil {
    //	log.Printf("epa")
    //}
  }


}

//instead of redrawing at every event, redraws every 300 microseconds
func refresh() {
  for {
    app.Draw()
    time.Sleep(200 * time.Microsecond) 
  }
}

func main() {

  go server.Start(c)
  go read_events()
  go refresh()


  eventList.
  AddItem("0 List item 1", "", 0, nil).
  AddItem("1 List item 2", "", 0, nil)

  for i := 1; i <= 100; i++ {
    
    eventList.AddItem(fmt.Sprintf("List item %d", i), "", 0, nil)
  }

  eventList.ShowSecondaryText(false)

  eventView.
  SetDynamicColors(true).
  SetRegions(true).
  SetWordWrap(true).
  SetChangedFunc(func() {
    app.Draw()
  })

    

  var obj map[string]interface{}
  json.Unmarshal([]byte("{ \"hola\": \"hola\"}"), &obj)

  // Make a custom formatter with indent set
  f := colorjson.NewFormatter()
  f.Indent = 4

  // Marshall the Colorized JSON
  s, _ := f.Marshal(obj)

  fmt.Fprintf(eventView, "%s", tview.TranslateANSI(string(s)))

  eventList.SetBorder(true)
  eventView.SetBorder(true)

  title :=  tview.NewTextView().
  SetDynamicColors(true).
  SetText(" Logshark ")

  footer :=  tview.NewTextView().
  SetDynamicColors(true).
  SetText(" [blue]R[white]efresh [blue]S[white]ettings ").
  SetChangedFunc(func() {
    app.Draw()
  })

  layout := tview.NewFlex().
  AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
    AddItem(title, 1, 1, false).
    AddItem(eventList, 0, 1, false).
    AddItem(eventView, 0, 3, false).
    AddItem(footer, 1, 1, false), 0, 2, false)


  if err := app.SetRoot(layout, true).SetFocus(eventList).EnableMouse(true).Run(); err != nil {
    panic(err)
  }



}