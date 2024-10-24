package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"net"
)

// connection to clichat
var conn net.Conn

// handler: server incoming connection
func handleIncomingMessage(conn net.Conn, message *tview.TextView) {
	defer conn.Close()

	// read the server message
	for {
		byte := make([]byte, 2048)
		_, err := conn.Read(byte)
		if err != nil {
			log.Println("Error reading incoming data: ", err)
			return
		}
		prev_message := message.GetText(true)
		message.SetText(prev_message + ">> " + string(byte) + "\n")
	}
}

// connect to clichat server
func connectToServer(welcomeBox *tview.TextView, message *tview.TextView) {
	prev_msg := welcomeBox.GetText(true)
	welcomeBox.SetText(prev_msg + ">> Starting Connection to CLICHAT\n")

	var err error
	conn, err = net.Dial("tcp", "localhost:80")
	if err != nil {
		log.Fatalf("Error connecting to CLICHAT server: %v", err)
	}

	prev_msg = welcomeBox.GetText(true)
	welcomeBox.SetText(prev_msg + ">> Connected to CLICHAT server\n")

	go handleIncomingMessage(conn, message)
}

// handle message rela
func messageRelay(conn net.Conn, msg string) {
	data := []byte(msg)
	_, err := conn.Write(data)
	if err != nil {
		log.Print("Error sending data to server", err)
	}
}

func main() {

	// init application
	app := tview.NewApplication()

	// application scaffolds
	// message view and input view
	right_view := tview.NewFlex().SetDirection(tview.FlexRow)
	right_view.SetBorder(true).SetTitle("CLICHAT").SetTitleAlign(tview.AlignLeft)
	left_view := tview.NewTextView()
	left_view.SetBorder(true).SetTitle("Message").SetTitleAlign(tview.AlignLeft)

	// main scaffold
	root := tview.NewFlex().SetDirection(tview.FlexColumn)

	// set up message box
	message_box := tview.NewInputField().
		SetLabel("Input Message:  ").
		SetFieldWidth(30).
		SetPlaceholder(" Enter message here...").
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorDarkCyan).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetPlaceholderTextColor(tcell.ColorYellow)
	message_box.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msg := message_box.GetText()
			messageRelay(conn, msg)
			message_box.SetText("")
		}
	})

	welcome_box := tview.NewTextView()
	connectToServer(welcome_box, left_view)

	// set up items designated view
	right_view.AddItem(welcome_box, 0, 1, false).
		AddItem(message_box, 0, 8, true)

	// add items to root
	root.AddItem(left_view, 0, 1, false).AddItem(right_view, 0, 2, true)

	// start application
	if err := app.SetRoot(root, true).Run(); err != nil {
		log.Fatal("Error starting application: ", err)
	}

}
