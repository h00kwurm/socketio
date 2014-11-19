# SocketIO .9

This is a super minimal implementation of my SocketIO needs. If you do find yourself wanting to use it, go right ahead. Be forewarned! It is really underdeveloped. Want to submit a pull request and help out others stuck in a shit mode of point-nine-age?

Yeah, I know there are [two](https://github.com/googollee/go-socket.io) [other](https://github.com/oguzbilgic/socketio) options for socketio. Both are backed by websockets provided by [code.google.com](https://code.google.com/p/go/). This is backed by [gorilla websockets](https://github.com/gorilla/websocket). Thank you [gorilla](https://github.com/gorilla).

If it wasn't clear by now, this only supports websockets. Maybe you're thinking to yourself, why socketio if it always uses just websockets. Because reasons. That's why.

## Example: 

    package main

    import (
      "fmt"
      "github.com/h00kwurm/socketio"
    )

    const remoteServer = "http://127.0.0.1:8088"

    func onConnect(output chan socketio.Message) {
      output <- socketio.CreateMessageEvent(`{"msg":"test message"}`)
    }

    func main() {

      socket := socketio.SocketIO{
        OnConnect: onConnect,
      }

      err := socketio.ConnectToSocket(remoteServer, &socket)
      if err != nil {
        fmt.Println(err)
        return
      }

    }



## Licensing
    The MIT License (MIT)

    Copyright (c) 2014 Aditya Natraj

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.