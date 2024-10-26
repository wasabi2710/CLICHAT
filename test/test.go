package main

import "github.com/rivo/tview"

func main() {

	app := tview.NewApplication()
	drop := tview.NewForm().AddDropDown("", []string{"client 1", "client 2", "client 3"}, 1, selected)

	if err := app.SetRoot(drop, true).Run(); err != nil {
		panic(err)
	}

}

func selected(option string, optionIndex int) {

}
