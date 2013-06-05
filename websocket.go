package main

import (
    "code.google.com/p/go.net/websocket"
    l4g "code.google.com/p/log4go"
//    "fmt"
)

var (
  ActiveClients = make(map[ClientConn]int)
)

type ClientConn struct {
  websocket *websocket.Conn
  clientIP string
}
type WebsocketCommand struct {
  Type string
  Callback_Id int
}

func publish(msg interface{}){
  var err error
  for cs, _ := range ActiveClients {
		if err = websocket.JSON.Send(cs.websocket, msg); err != nil {
			l4g.Error("Could not send message to ", cs.clientIP, err.Error())
		}
	}

  l4g.Info("sent: %#v", msg)
}

func wsHandler(ws *websocket.Conn){
  l4g.Info("ws connection: %v", ws.Config())

  var err error

  defer func(){
    if err = ws.Close(); err != nil {
      l4g.Info("Websocket couldn't be closed",err.Error())
    }
  }()

  client := ws.Request().RemoteAddr
  l4g.Info("Client connected: ", client)

  conn := ClientConn{ws, client}
  ActiveClients[conn] = 0

  l4g.Info("%v clients connected.", len(ActiveClients))

  for {
    var msg WebsocketCommand
    // Receive receives a text message serialized T as JSON
    if err = websocket.JSON.Receive(ws, &msg); err != nil {
      l4g.Info("Error receiving: %v", err)
      delete(ActiveClients, conn)
      return
    }

    l4g.Info("recv from %v: %#v",client, msg)
    // TODO: put this request and conn in the queue that the server looks at
    //  so it can get to it whenever.
    // OR: just do a service call and send the results.
    //     less complicated that way. The less you do, the faster things are.
    
    response := requestHandler(&msg)
    err = websocket.JSON.Send(ws, response)
    if err != nil {
      l4g.Info("Error sending: %#v", err)
    }

    l4g.Info("sent: %#v", response)

  }

}
