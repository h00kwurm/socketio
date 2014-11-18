package socketio

import "testing"

func TestIncrementMessageIndex(t *testing.T) {
	if currentIndex != 1 {
		t.Errorf("current index not as expected: ex: %v ac: %v\n", 1, currentIndex)
	}
	for i := 0; i < 129; i++ {
		incrementMessageIndex()
	}
	if currentIndex != 1 {
		t.Errorf("current index not as expected: ex: %v ac: %v\n", 1, currentIndex)
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
