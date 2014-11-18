// The MIT License (MIT)

// Copyright (c) 2013 Oguz Bilgic

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package socketio

import (
	"testing"
)

func testHandshakeUrls(t *testing.T, rawUrls []string, expected string) {
	for _, raw := range rawUrls {
		u, err := newURLParser(raw)
		if err != nil {
			t.Errorf("NewUrl error:  %s", err)
		}
		if u.handshake() != expected {
			t.Errorf("Wrong handshake formatted url, expected: %s, actual: %s", expected, u.handshake())
		}
	}
}

func testWebsocketUrls(t *testing.T, rawUrls []string, expected string) {
	for _, raw := range rawUrls {
		u, err := newURLParser(raw)
		if err != nil {
			t.Errorf("NewUrl error:  %s", err)
		}
		ws := u.websocket("session_id")
		if ws != expected {
			t.Errorf("Wrong websocket formatted url, expected: %s, actual: %s", expected, ws)
		}
	}
}

func TestHandshakeUrl(t *testing.T) {
	testHandshakeUrls(t,
		[]string{"server.com", "http://server.com"},
		"http://server.com/socket.io/1")

	testHandshakeUrls(t,
		[]string{"server.com/path", "http://server.com/path"},
		"http://server.com/path/socket.io/1")

	testHandshakeUrls(t,
		[]string{"https://server.com"},
		"https://server.com/socket.io/1")
}

func TestWebsocketUrl(t *testing.T) {
	testWebsocketUrls(t,
		[]string{"server.com", "http://server.com"},
		"ws://server.com/socket.io/1/websocket/session_id")

	testWebsocketUrls(t,
		[]string{"server.com/path", "http://server.com/path"},
		"ws://server.com/path/socket.io/1/websocket/session_id")

	testWebsocketUrls(t,
		[]string{"https://server.com"},
		"wss://server.com/socket.io/1/websocket/session_id")
}
