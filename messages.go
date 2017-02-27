package socketio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
)

type Message struct {
	Type int
	Event string
	ID   int
	Body []byte
	Bodyv1 []byte
	Version float64
	Ack  func(message []byte, output chan Message)
}

var currentIndex = 1

func incrementMessageIndex() {
	if currentIndex == 128 {
		currentIndex = 1
	} else {
		currentIndex++
	}
}

type Event struct {
	Event string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}

func CreateMessageEvent(name, message string, ack func(message []byte, output chan Message), version float64) Message {
	var temp json.RawMessage
	json.Unmarshal([]byte(message), &temp)
	tempArray := []json.RawMessage{temp}
	messageBody := Event{}

	if name == "" {
		messageBody = Event{
			Event: "message",
			Args: tempArray,
		}
	} else {
		messageBody = Event{
			Event: name,
			Args: tempArray,
		}
	}

	messageEvent := Message{
		Type: 5,
		Event: name,
		ID:   currentIndex,
		Ack:  ack,
		Version: version,
	}

	if version == 1 {
		messageEvent.Bodyv1 = []byte(message)
	}

	incrementMessageIndex()

	tempMessage, err := json.Marshal(messageBody)
	if err != nil {
		fmt.Println("error on marshal: ", err)
		return Message{}
	}

	messageEvent.Body = tempMessage

	return messageEvent

}

func CreateMessageHeartbeat() Message {
	message := Message{
		Type: 2,
	}
	return message
}

func (message Message) PrintMessage() string {
	if message.Version == 1 {
		return `42["` + message.Event + `",` + string(message.Bodyv1) + "]"
	} else {
		switch message.Type {
		case 2:
			return "2::"
		case 5:
			return "5:" + strconv.Itoa(message.ID) + "+::" + string(message.Body)
		default:
			return ""
		}
	}
}

func parseEvent(buffer []byte) (string, []byte) {
	var event Event
	index := bytes.Index(buffer, []byte("{"))
	json.Unmarshal(buffer[index:], &event)
	return event.Event, event.Args[0]
}

func parseEventv1(buffer []byte) (string, []byte) {
	msg := strings.TrimLeft(string(buffer), "42")
	var temp []interface{}
	json.Unmarshal([]byte(msg), &temp)
	data, _ := json.Marshal(temp[1])
	eventName := temp[0].(string)
	return eventName, data
}

func parseMessage(buffer []byte) []byte {
	splitChunks := bytes.Split(buffer, []byte(":"))
	if len(splitChunks) < 4 {
		return []byte("")
	}
	return splitChunks[3]
}

func parseAck(buffer []byte) (int, []byte) {
	if len(buffer) < 5 {
		return 0, []byte("")
	}
	id, _ := strconv.Atoi(string(buffer[4]))

	if len(buffer) > 5 {
		if uint8(buffer[5]) != 43 {
			return id, []byte("")
		} else {
			index := bytes.Index(buffer, []byte("{"))
			if index != -1 {
				lastIndex := len(buffer) - 1
				return id, buffer[index:lastIndex]
			} else {
				return id, []byte("")
			}
		}
	} else {
		return id, []byte("")
	}

}
