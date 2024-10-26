package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MessageType represents the type of message being sent
type MessageType int

const (
	WelcomeMessage MessageType = iota
	ClientListMessage
	ChatMessage
)

// Message represents a message with a type and payload
type Message struct {
	Type      MessageType `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp string      `json:"timestamp"`
}

// connection to clichat
var conn net.Conn

// handler: server incoming connection
func handleIncomingMessage(conn net.Conn, message *tview.TextView, clientList *tview.Form) {
	defer conn.Close()

	for {
		var length int32
		err := binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			log.Println("Error reading message length: ", err)
			return
		}

		buffer := make([]byte, length)
		_, err = conn.Read(buffer)
		if err != nil {
			log.Println("Error reading incoming data: ", err)
			return
		}

		var msg Message
		err = json.Unmarshal(buffer, &msg)
		if err != nil {
			log.Println("Error unmarshalling message: ", err)
			continue
		}

		prevMessage := message.GetText(true)
		switch msg.Type {
		case WelcomeMessage:
			message.SetText(prevMessage + "### " + msg.Payload.(string) + "\n")
		case ChatMessage:
			message.SetText(prevMessage + msg.Timestamp + " >> " + msg.Payload.(string) + "\n")
		case ClientListMessage:
			//var _clients []string
			//clientAddrs := msg.Payload.([]interface{})
			//for _, addr := range clientAddrs {
			//	_clients = append(_clients, addr.(string))
			//}
		}
	}
}

// connect to clichat server
func connectToServer(welcomeBox *tview.TextView, message *tview.TextView, clientList *tview.Form) {
	prevMsg := welcomeBox.GetText(true)
	welcomeBox.SetText(prevMsg + "### Starting Connection to CLICHAT\n")

	var err error
	conn, err = net.Dial("tcp", "localhost:80")
	if err != nil {
		log.Fatalf("Error connecting to CLICHAT server: %v", err)
	}

	prevMsg = welcomeBox.GetText(true)
	welcomeBox.SetText(prevMsg + "### Connected to CLICHAT server\n")

	go handleIncomingMessage(conn, message, clientList)
}

// handle message relay
func messageRelay(conn net.Conn, msg string) {
	message := Message{
		Type:    ChatMessage,
		Payload: msg,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Print("Error marshalling message: ", err)
		return
	}

	length := int32(len(data))
	err = binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		log.Print("Error writing message length: ", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Print("Error sending data to server: ", err)
	}
}

func main() {
	app := tview.NewApplication()

	rightView := tview.NewFlex().SetDirection(tview.FlexRow)

	// two more boxes for right view
	welcomeView := tview.NewFlex().SetDirection(tview.FlexRow)
	welcomeView.SetBorder(true).SetTitle("CLICHAT").SetTitleAlign(tview.AlignLeft)

	leftView := tview.NewTextView()
	leftView.SetBorder(true).SetTitle("Message").SetTitleAlign(tview.AlignLeft)

	clientList := tview.NewForm()
	clientList.SetBorder(true).SetTitle("Clients").SetTitleAlign(tview.AlignLeft)

	root := tview.NewFlex().SetDirection(tview.FlexColumn)

	messageBox := tview.NewInputField().
		SetLabel("Input Message:  ").
		SetFieldWidth(30).
		SetPlaceholder(" Enter message here...").
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorDarkCyan).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetPlaceholderTextColor(tcell.ColorYellow)
	messageBox.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msg := messageBox.GetText()
			messageRelay(conn, msg)
			messageBox.SetText("")
		}
	})

	welcomeBox := tview.NewTextView()
	connectToServer(welcomeBox, leftView, clientList)

	welcomeView.AddItem(welcomeBox, 0, 1, false).
		AddItem(messageBox, 0, 6, false)

	rightView.AddItem(welcomeView, 0, 1, false).
		AddItem(clientList, 0, 1, true)

	root.AddItem(leftView, 0, 1, false).
		AddItem(rightView, 0, 2, true)

	if err := app.SetRoot(root, true).Run(); err != nil {
		log.Fatal("Error starting application: ", err)
	}
}
