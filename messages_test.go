package socketio

import (
	"testing"
)

func TestIncrementMessageIndex(t *testing.T) {
	if currentIndex != 1 {
		t.Errorf("current index not as expected: ex: %v ac: %v\n", 1, currentIndex)
	}
	for i := 0; i < 64; i++ {
		incrementMessageIndex()
	}
	if currentIndex != 65 {
		t.Errorf("current index not as expected: ex: %v ac: %v\n", 65, currentIndex)
	}
	for i := 0; i < 64; i++ {
		incrementMessageIndex()
	}
	if currentIndex != 1 {
		t.Errorf("current index not as expected: ex: %v ac: %v\n", 1, currentIndex)
	}
	incrementMessageIndex()
	if currentIndex != 2 {
		t.Errorf("current index not as expected: ex: %v ac: %v\n", 2, currentIndex)
	}
}

func TestCreateMessageEvent(t *testing.T) {
	const content = `{"msg":"fuck you","body":[{"this":"test"}]}`
	msg := CreateMessageEvent(content)

	if msg.Type != 5 {
		t.Errorf("[test] message:event. failed. invalid type: %v\n", msg.Type)
	}

	if msg.ID == 0 {
		t.Errorf("[test] message:event. failed. invalid id: %v\n", msg.ID)
	}

	if currentIndex != (msg.ID + 1) {
		t.Errorf("[test] message:event. didn't increment cIndex. index: %v, id: %v\n", currentIndex, msg.ID)
	}

	name, outputContent := parseEvent([]byte(msg.PrintMessage()))
	if name != "message" {
		t.Errorf("[test] createMessgeEvent: bad name. got %v\n", name)
	}

	for i := 0; i < len(content); i++ {
		if outputContent[i] != content[i] {
			t.Errorf("[test] createMessageEvent: bad content. got: %v\n expected: %v\n", string(outputContent), content)
			break
		}
	}

}

func TestHeartbeat(t *testing.T) {
	msg := CreateMessageHeartbeat()

	if msg.Type != 2 {
		t.Errorf("[test] heartbeat. failed. invalid type: %v\n", msg.Type)
	}

	if msg.PrintMessage() != "2::" {
		t.Errorf("[test] heartbeat. failed. not ok. ex: %v, ac: %v\n", "2::", msg.PrintMessage())
	}
}

func TestParseEvent(t *testing.T) {
	const testEvent = `5:1+::{"name":"message","args":[{"msg":"WebsocketAuthenticate","seq":1,"datatype":"WSAuthenticateType3","dst":"A8:77:6F:00:27:22","body":[{"ticket":"ROqdMmUdHIa/KwpVHR/YOq78iuKGSKhYg1cmMJFW+pxZwosrRx/5XETLZxRLt8q4","version":"7.5.2"}]}]}`
	const testContent = `{"msg":"WebsocketAuthenticate","seq":1,"datatype":"WSAuthenticateType3","dst":"A8:77:6F:00:27:22","body":[{"ticket":"ROqdMmUdHIa/KwpVHR/YOq78iuKGSKhYg1cmMJFW+pxZwosrRx/5XETLZxRLt8q4","version":"7.5.2"}]}`

	name, content := parseEvent([]byte(testEvent))

	if name != "message" {
		t.Errorf("[test] parseEvent: bad name. got %v\n", name)
	}

	for i := 0; i < len(content); i++ {
		if content[i] != testContent[i] {
			t.Errorf("[test] parseEvent: bad content. got %v\n", string(content))
			break
		}
	}
}

func TestParseMessage(t *testing.T) {
	const testGoodMessage = `3:1::blabla`
	const testBadMessage = `3:1:blabla`
	const testEmptyMessage = `3:::`

	if len(parseMessage([]byte(testEmptyMessage))) != 0 {
		t.Error("[test] empty message. failed. not empty.")
	}

	if len(parseMessage([]byte(testBadMessage))) != 0 {
		t.Error("[test] bad message. failed. not empty.")
	}

	output := parseMessage([]byte(testGoodMessage))
	if len(output) == 0 {
		t.Error("[test] good message. failed. is empty")
	}

	if string(output) != "blabla" {
		t.Errorf("[test] good message not as expected. ex: %v, ac: %v\n", "blabla", string(output))
	}

}
