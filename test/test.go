package main

import (
	"log"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	list := tview.NewList().
		AddItem("List item 1", "Some explanatory text", 'a', nil).
		AddItem("List item 2", "Some explanatory text", 'b', nil).
		AddItem("List item 3", "Some explanatory text", 'c', nil).
		AddItem("List item 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	// Add the SetSelectedFunc to handle item selection
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		log.Printf("Selected item: %s", mainText)
		// Add your logic here for when an item is selected
	})

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
