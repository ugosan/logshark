// Demo code for the List primitive.
package main

import (
	"github.com/ugosan/logshark/cmd/server"
	"github.com/rivo/tview"
	"fmt"
	"encoding/json"
	"github.com/TylerBrock/colorjson"
)

func main() {


	go server.Start()
	

	app := tview.NewApplication()
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	menu := tview.NewList().
	AddItem("List item 1", "", 0, nil).
	AddItem("List item 2", "Some explanatory text", 0, nil).
	AddItem("List item 3", "Some explanatory text", 0, nil).
	AddItem("List item 4", "Some explanatory text", 0, nil)

	for i := 1; i <= 100; i++ {
		
		menu.AddItem(fmt.Sprintf("List item %d", i), "", 0, nil)
	}

	menu.ShowSecondaryText(false)

	main := tview.NewTextView().
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

	fmt.Fprintf(main, "%s", string(s))


	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, true).
		AddItem(main, 1, 1, 1, 2, 0, 100, true)

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}