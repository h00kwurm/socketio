package socketio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type SocketIO struct {
	Context    *Session
	Connection *websocket.Conn

	InputChannel  chan string
	OutputChannel chan string

	OnConnect    func(output chan string)
	OnDisconnect func(output chan string)
	OnMessage    func(message []byte, output chan string)
	OnJSON       func(message []byte, output chan string)
	OnEvent      map[string]func(message []byte, output chan string)
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

	socket.OutputChannel = make(chan string)
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
		case _, outgoing_state := <-socket.OutputChannel:
			if !outgoing_state {
				return errors.New("output channel closed")
			}
		}
	}

	return err
}

func processBus(input chan string, output chan string) error {

	for {
		select {
		case incoming, incoming_state := <-input:
			if !incoming_state {
				fmt.Println("input channel is broken")
				return errors.New("input channel is broken")
			}
			fmt.Println(string(incoming))
		case _, outgoing_state := <-output:
			if !outgoing_state {
				return errors.New("output channel closed")
			}
		}
	}

}

func (socket *SocketIO) readInput() {
	for {
		msgType, buffer, err := socket.Connection.ReadMessage()
		if err != nil {
			fmt.Println("error!: ", err)
			break
		}

		switch uint8(buffer[0]) {
		case 48: //0:
			fmt.Println("socket closed!")
			break
		case 49: //1:
			fmt.Println("socket opened")
			if socket.OnConnect != nil {
				socket.OnConnect(socket.OutputChannel)
			}
		case 50: //2:
			fmt.Println("heartbeat received")
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
		}

		if msgType == 1 {
			fmt.Println("message :", buffer)
			socket.InputChannel <- string(buffer)
		}
	}
}

type Event struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}

func parseEvent(buffer []byte) (string, []byte) {
	var event Event
	json.Unmarshal([]byte(testData), &event)
	return event.Name, event.Args[0]
}

func parseMessage(buffer []byte) []byte {
	splitChunks := bytes.Split(buffer, []byte(":"))
	return splitChunks[3]
}

func buildUrl(url string, endpoint string) string {
	if strings.Contains(url, "http") {
		return strings.Replace(url, "http", "ws", 1) + "/socket.io/1/websocket/" + endpoint
	} else if strings.Contains(url, "https") {
		return strings.Replace(url, "https", "wss", 1) + "/socket.io/1/websocket/" + endpoint
	}
	return url
}
