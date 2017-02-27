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
	"fmt"
	"net/url"
	"strings"
)

// Parse raw url string and make valid handshake or websockets socket.io url.
type urlParser struct {
	raw    string
	parsed *url.URL
}

func newURLParser(raw string) (*urlParser, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme == "" {
		parsed.Scheme = "http"
	}
	return &urlParser{raw: raw, parsed: parsed}, nil
}

func (u *urlParser) handshake(version float64) string {
	if version == 1 {
		return fmt.Sprintf("%s/socket.io/?transport=polling&b64=1", u.parsed.String()) 
	} else {
		return fmt.Sprintf("%s/socket.io/1", u.parsed.String())
	}
}

func (u *urlParser) websocket(sessionId string, version float64) string {
	var host string
	if u.parsed.Scheme == "https" {
		host = strings.Replace(u.parsed.String(), "https://", "wss://", 1)
	} else {
		host = strings.Replace(u.parsed.String(), "http://", "ws://", 1)
	}
	if version == 1 {
		return fmt.Sprintf("%s/socket.io/?transport=websocket&sid=%s", host, sessionId)
	} else {
		return fmt.Sprintf("%s/socket.io/1/websocket/%s", host, sessionId)
	}
}
