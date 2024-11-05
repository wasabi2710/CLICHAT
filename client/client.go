package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
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
	SelectedClient
)

// Message represents a message with a type and payload
type Message struct {
	Type      MessageType `json:"type"`
	Sender    string      `json:"sender"`
	Hostaddr  string      `json:"hostaddr"`
	Relay     string      `json:"relay"`
	Payload   interface{} `json:"payload"`
	Timestamp string      `json:"timestamp"`
}

// change tab focus
type Focus struct {
	view     map[int]*tview.Box
	selected int
}

// connection to clichat
var conn net.Conn
var host string

// available clients list
var _clients []string

// current client
var curClient string
var curName string

// handler: clients selection
func clientSelect(index int, mainText string, secondaryText string, shortcut rune) {
	curClient = mainText
	// tell the server that i need to relay my messages this client currently
	//selectedClient(conn, curClient)
}

// handler: server incoming connection
func handleIncomingMessage(conn net.Conn, message *tview.List, clientList *tview.List) {
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
		switch msg.Type {
		case WelcomeMessage:
			curStr := fmt.Sprintf("{ %s }", msg.Payload.(string))
			message.AddItem(string(curStr), "", '\u0000', nil)
		case ChatMessage:
			if msg.Sender == conn.LocalAddr().String() {
				curStr := fmt.Sprintf("(YOU) -> \033[1;34m%s\033[0m", msg.Payload.(string))
				message.AddItem(curStr, msg.Timestamp, '\u0000', nil)
			} else {
				curStr := fmt.Sprintf("(%s) -> %s", msg.Hostaddr, msg.Payload.(string))
				message.AddItem(curStr, msg.Timestamp, '\u0000', nil)
			}
		case ClientListMessage:
			clientAddrs := msg.Payload.([]interface{})
			for _, addr := range clientAddrs {
				clientMap := addr.(map[string]interface{})
				for hostname, address := range clientMap {
					addrStr, err := address.(string)
					if !err {
						log.Print("Error: address is not a string")
						continue
					}
					if addrStr != conn.LocalAddr().String() {
						_clients = append(_clients, addrStr)
						clientList.AddItem(addrStr, hostname, '>', nil)
					}
				}
			}
			clientList.SetSelectedFunc(clientSelect).SetSelectedTextColor(tcell.ColorLawnGreen)
		}
	}
}

// connect to clichat server
func connectToServer(welcomeBox *tview.TextView, message *tview.List, clientList *tview.List) {
	prevMsg := welcomeBox.GetText(true)
	welcomeBox.SetText(prevMsg + "### Starting Connection to CLICHAT\n")
	var err error
	conn, err = net.Dial("tcp", "10.0.2.15:9999")
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
		Relay:   curClient,
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

// hanlde client selection relay
func selectedClient(conn net.Conn, client string) {
	// define payload
	thisClient := Message{
		Type:    SelectedClient,
		Payload: client,
	}
	// serialize payload
	data, err := json.Marshal(thisClient)
	if err != nil {
		log.Print("Error marshalling client: ", err)
		return
	}
	// get length
	length := int32(len(data))
	err = binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		log.Print("Error writing message length: ", err)
		return
	}
	// write data to server
	_, err = conn.Write(data)
	if err != nil {
		log.Print("Error writing data to server: ", err)
		return
	}
}

func main() {
	app := tview.NewApplication()
	rightView := tview.NewFlex().SetDirection(tview.FlexRow)
	// two more boxes for right view
	welcomeView := tview.NewFlex().SetDirection(tview.FlexRow)
	welcomeView.SetBorder(true).SetTitle("CLICHAT").SetTitleAlign(tview.AlignLeft)
	leftView := tview.NewList()
	leftView.SetBorder(true).SetTitle("Message").SetTitleAlign(tview.AlignLeft)
	clientList := tview.NewList()
	clientList.SetBorder(true).SetTitle("Clients").SetTitleAlign(tview.AlignLeft)
	root := tview.NewFlex().SetDirection(tview.FlexColumn)
	messageBox := tview.NewInputField().SetLabel("Input Message:  ").SetFieldWidth(30).SetPlaceholder(" Enter message here...").SetFieldTextColor(tcell.ColorWhite).SetLabelColor(tcell.ColorDarkCyan).SetFieldBackgroundColor(tcell.ColorBlack).SetPlaceholderTextColor(tcell.ColorYellow)
	messageBox.SetBorder(true).SetTitle("Message Input").SetTitleAlign(tview.AlignLeft)
	messageBox.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msg := messageBox.GetText()
			messageRelay(conn, msg)
			messageBox.SetText("")
		}
	})

	tab := &Focus{
		view: map[int]*tview.Box{
			1: messageBox.Box,
			2: clientList.Box,
			3: leftView.Box,
		},
		selected: 1,
	}

	// switch focus here
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			tab.selected++
			if tab.selected > 3 {
				tab.selected = 1
			}
			app.SetFocus(tab.view[tab.selected])
		}
		return event
	})

	welcomeBox := tview.NewTextView()
	connectToServer(welcomeBox, leftView, clientList)
	welcomeView.AddItem(welcomeBox, 0, 1, false)
	rightView.AddItem(welcomeView, 0, 1, true).AddItem(clientList, 0, 1, false).AddItem(messageBox, 0, 1, true)
	root.AddItem(leftView, 0, 1, false).AddItem(rightView, 0, 1, true)
	if err := app.SetRoot(root, true).Run(); err != nil {
		log.Fatal("Error starting application: ", err)
	}
}
