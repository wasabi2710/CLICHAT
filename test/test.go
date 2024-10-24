package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type Client struct {
	name string
}

var selected string

func main() {

	app := tview.NewApplication()

	clients := []Client{
		{name: "client 1"},
		{name: "client 2"},
		{name: "client 3"},
	}

	clientNames := make([]string, len(clients))
	for i, client := range clients {
		clientNames[i] = client.name
	}

	form := tview.NewForm().
		AddDropDown("Clients", clientNames, 0, choose).
		AddButton("Call", func() {
			fmt.Println("you chose ", selected)
		})

	if err := app.SetRoot(form, true).Run(); err != nil {
		panic(err)
	}
}

func choose(option string, optionIndex int) {
	selected = option
}

