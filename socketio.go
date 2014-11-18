package socketio

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type SocketIO struct {
	Context    *Session
	Connection *websocket.Conn

	InputChannel  chan string
	OutputChannel chan Message

	callbacks map[int]func(message []byte)

	OnConnect    func(output chan Message)
	OnDisconnect func(output chan Message)
	OnMessage    func(message []byte, output chan Message)
	OnJSON       func(message []byte, output chan Message)
	OnEvent      map[string]func(message []byte, output chan Message)
	OnError      func()
}

func ConnectToSocket(urlString string, socket *SocketIO) error {

	var err error

	socket.Context, err = NewSession(urlString)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var connector = websocket.Dialer{
		HandshakeTimeout: (*socket.Context).HeartbeatTimeout,
		Subprotocols:     []string{"websocket"},
	}

	connectionUrl := buildUrl(urlString, (*socket.Context).ID)

	socket.Connection, _, err = connector.Dial(connectionUrl, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer socket.Connection.Close()

	socket.InputChannel = make(chan string)
	defer close(socket.InputChannel)

	socket.OutputChannel = make(chan Message)
	defer close(socket.OutputChannel)

	go socket.readInput()
	// go processBus(socket.InputChannel, socket.OutputChannel)

	for {
		select {
		case incoming, incoming_state := <-socket.InputChannel:
			if !incoming_state {
				fmt.Println("input channel is broken")
				return errors.New("input channel is broken")
			}
			fmt.Println(string(incoming))
		case outgoing, outgoing_state := <-socket.OutputChannel:
			if !outgoing_state {
				return errors.New("output channel closed")
			}
			item := outgoing.PrintMessage()
			fmt.Println("sending --> ", item)
			if err := socket.Connection.WriteMessage(1, []byte(item)); err != nil {
				fmt.Println(err)
				return errors.New("io corrupted. can't continue")
			}
		}
	}

	return err
}

// I wanted to use this as :> go processBus as you can see around
// the late 50s lines but it isn't working right now. whatever
// i'll figure it out
// func (socket *SocketIO) processBus() {

// 	for {
// 		select {
// 		case incoming, incoming_state := <-socket.InputChannel:
// 			if !incoming_state {
// 				fmt.Println("input channel is broken")
// 				return
// 			}
// 			fmt.Println(string(incoming))
// 		case outgoing, outgoing_state := <-socket.OutputChannel:
// 			if !outgoing_state {
// 				fmt.Println("output channel closed")
// 				return
// 			}
// 			if err := socket.Connection.WriteMessage(1, []byte(outgoing.PrintMessage())); err != nil {
// 				fmt.Println("io corrupted, can't continue: ", err)
// 				return
// 			}
// 		}
// 	}
// }

func (socket *SocketIO) readInput() {
	for {
		msgType, buffer, err := socket.Connection.ReadMessage()
		if err != nil {
			fmt.Println("error!: ", err)
			break
		}
		fmt.Println("received-->", string(buffer))

		switch uint8(buffer[0]) {
		case 48: //0:
			if socket.OnDisconnect != nil {
				socket.OnDisconnect(socket.OutputChannel)
			}
			break
		case 49: //1:
			if socket.OnConnect != nil {
				socket.OnConnect(socket.OutputChannel)
			}
		case 50: //2:
			socket.OutputChannel <- CreateMessageHeartbeat()
		case 51: //3:
			if socket.OnMessage != nil {
				message := parseMessage(buffer)
				socket.OnMessage(message, socket.OutputChannel)
			}
		case 52: //4:
			if socket.OnJSON != nil {
				message := parseMessage(buffer)
				socket.OnJSON(message, socket.OutputChannel)
			}
		case 53: //5:
			if socket.OnEvent != nil {
				eventName, eventMessage := parseEvent(buffer)
				if eventFunction := socket.OnEvent[eventName]; eventFunction != nil {
					eventFunction(eventMessage, socket.OutputChannel)
				}
			}
		case 54: //6:

		case 55: //7:
			if socket.OnError != nil {
				socket.OnError()
			}
			break
		}

	}
}

func buildUrl(url string, endpoint string) string {
	if strings.Contains(url, "http") {
		return strings.Replace(url, "http", "ws", 1) + "/socket.io/1/websocket/" + endpoint
	} else if strings.Contains(url, "https") {
		return strings.Replace(url, "https", "wss", 1) + "/socket.io/1/websocket/" + endpoint
	}
	return url
}
