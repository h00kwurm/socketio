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

	callbacks map[int]func(message []byte, output chan Message)

	OnConnect    func(output chan Message)
	OnDisconnect func(output chan Message)
	OnMessage    func(message []byte, output chan Message)
	OnJSON       func(message []byte, output chan Message)
	OnAck        func(message []byte, output chan Message)
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

	socket.callbacks = make(map[int]func(message []byte, output chan Message))

	socket.InputChannel = make(chan string)
	defer close(socket.InputChannel)

	socket.OutputChannel = make(chan Message)
	defer close(socket.OutputChannel)

	go socket.readInput()
	// go socket.processBus()

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
			if outgoing.Type == 5 && outgoing.Ack != nil {
				socket.callbacks[outgoing.ID] = outgoing.Ack
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
// 			if outgoing.Type == 5 && outgoing.Ack != nil {
// 				socket.callbacks[outgoing.ID] = outgoing.Ack
// 			}
// 			item := outgoing.PrintMessage()
// 			fmt.Println("sending --> ", item)
// 			if err := socket.Connection.WriteMessage(1, []byte(item)); err != nil {
// 				fmt.Println(err)
// 				fmt.Println("io corrupted. can't continue")
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

		if msgType == 1 {
			switch uint8(buffer[0]) {
			case 48: //0:
				if socket.OnDisconnect != nil {
					go socket.OnDisconnect(socket.OutputChannel)
				}
				break
			case 49: //1:
				if socket.OnConnect != nil {
					go socket.OnConnect(socket.OutputChannel)
				}
			case 50: //2:
				socket.OutputChannel <- CreateMessageHeartbeat()
			case 51: //3:
				if socket.OnMessage != nil {
					message := parseMessage(buffer)
					go socket.OnMessage(message, socket.OutputChannel)
				}
			case 52: //4:
				if socket.OnJSON != nil {
					message := parseMessage(buffer)
					go socket.OnJSON(message, socket.OutputChannel)
				}
			case 53: //5:
				if socket.OnEvent != nil {
					eventName, eventMessage := parseEvent(buffer)
					if socket.OnEvent != nil {
						if eventFunction := socket.OnEvent[eventName]; eventFunction != nil {
							go eventFunction(eventMessage, socket.OutputChannel)
						}
					}
				}
			case 54: //6:
				id, data := parseAck(buffer)
				function, exists := socket.callbacks[id]
				if exists {
					go function(data, socket.OutputChannel)
					delete(socket.callbacks, id)
				}
				if socket.OnAck != nil {
					go socket.OnAck(data, socket.OutputChannel)
				}
			case 55: //7:
				if socket.OnError != nil {
					go socket.OnError()
				}
				break
			}

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
