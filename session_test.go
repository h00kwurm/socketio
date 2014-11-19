package socketio

import (
	"testing"
	"time"
)

func TestSupportProtocol(t *testing.T) {
	session := Session{
		ID:                 "testID",
		HeartbeatTimeout:   time.Second,
		ConnectionTimeout:  time.Second,
		SupportedProtocols: []string{"websocket", "xhr-polling", "flashsocket"},
	}

	if !session.SupportProtocol("websocket") {
		t.Error("[test:supportprotocol] not matching protocols properly")
	}

	if session.SupportProtocol("jsonp-polling") {
		t.Error("[test:supportprotocol] not matching protocols properly")
	}

}
