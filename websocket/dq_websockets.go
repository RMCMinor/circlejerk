package dq_websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type WebsocketSubscriber struct {
	ws   *websocket.Conn
	msgs chan any
}

type WsServer struct {
	clients map[WebsocketSubscriber]struct{}
}

var wsServer = WsServer{
	clients: make(map[WebsocketSubscriber]struct{}),
}

func (server *WsServer) AddListener(w http.ResponseWriter, r *http.Request) {

	c, err := websocket.Accept(w, r, nil)

	if err != nil {
		http.Error(w, "Could not accept websocket connection: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer c.CloseNow()

	wsSubscriber := WebsocketSubscriber{
		msgs: make(chan any, 16),
		ws: c,
	}

	log.Printf("adding listerer\n")
	server.clients[wsSubscriber] = struct{}{}
	defer wsServer.RemoveListener(wsSubscriber)

	ctx := c.CloseRead(context.Background())

	for {
		select {
		case msg := <-wsSubscriber.msgs:
			err := wsjson.Write(context.Background(), wsSubscriber.ws, msg)

			if err != nil {
				return
			}
		case <-ctx.Done():
			log.Println("done")
			return
		}
	}
}

func (server *WsServer) RemoveListener(ws WebsocketSubscriber) {
	log.Printf("removing listerer\n")
	delete(server.clients, ws)
}

func WebsocketConnect(w http.ResponseWriter, r *http.Request) {
	wsServer.AddListener(w, r)
}

func SendWSMessage(msg any) {
	json, _ := json.Marshal(msg)
	log.Printf("%s\n", json)
	for client := range wsServer.clients {
		client.msgs <- msg
	}
}
