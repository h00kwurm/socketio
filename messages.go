package socketio

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	Type int
	ID   int
	Body []byte
}

type Event struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}

var currentIndex = 1

func incrementMessageIndex() {
	if currentIndex > 128 {
		currentIndex = 1
	} else {
		currentIndex++
	}
}

func CreateMessageEvent(message string) Message {

	var temp json.RawMessage
	json.Unmarshal([]byte(message), &temp)
	tempArray := []json.RawMessage{temp}

	messageBody := Event{
		Name: "message",
		Args: tempArray,
	}

	messageEvent := Message{
		Type: 5,
		ID:   currentIndex,
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

func (message Message) PrintMessage() string {
	if message.Type == 5 {
		return "5:" + strconv.Itoa(message.ID) + "+::" + string(message.Body)
	}
	return ""
}

func parseEvent(buffer []byte) (string, []byte) {
	var event Event
	json.Unmarshal([]byte(buffer), &event)
	return event.Name, event.Args[0]
}

func parseMessage(buffer []byte) []byte {
	splitChunks := bytes.Split(buffer, []byte(":"))
	if len(splitChunks) < 4 {
		return []byte("")
	}
	return splitChunks[3]
}
