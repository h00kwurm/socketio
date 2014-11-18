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
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	ID                 string
	HeartbeatTimeout   time.Duration
	ConnectionTimeout  time.Duration
	SupportedProtocols []string
}

// NewSession receives the configuration variables from the socket.io
// server.
func NewSession(url string) (*Session, error) {
	urlParser, err := newURLParser(url)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(urlParser.handshake())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()

	sessionVars := strings.Split(string(body), ":")
	if len(sessionVars) != 4 {
		return nil, errors.New("Session variables is not 4")
	}

	id := sessionVars[0]

	heartbeatTimeoutSec, _ := strconv.Atoi(sessionVars[1])
	connectionTimeoutSec, _ := strconv.Atoi(sessionVars[2])

	heartbeatTimeout := time.Duration(heartbeatTimeoutSec) * time.Second
	connectionTimeout := time.Duration(connectionTimeoutSec) * time.Second

	supportedProtocols := strings.Split(string(sessionVars[3]), ",")

	return &Session{id, heartbeatTimeout, connectionTimeout, supportedProtocols}, nil
}

// SupportProtocol checks if the given protocol is supported by the
// socket.io server.
func (session *Session) SupportProtocol(protocol string) bool {
	for _, supportedProtocol := range session.SupportedProtocols {
		if protocol == supportedProtocol {
			return true
		}
	}
	return false
}
