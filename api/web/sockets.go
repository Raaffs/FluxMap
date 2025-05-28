package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)


type Websocket struct{
	Upgrader websocket.Upgrader
	Clients  map[*websocket.Conn]bool
	mutex 	 sync.Mutex
}

func (ws *Websocket)Upgrade(r *http.Request){
	ws.Upgrader=websocket.Upgrader{
		CheckOrigin: func(r *http.Request)bool{return true},
	}	
}

func (ws *Websocket)AddClient(conn *websocket.Conn){
	ws.mutex.Lock()
	ws.Clients[conn]=true
	ws.mutex.Unlock()
}

func (ws *Websocket)RemoveClient(conn *websocket.Conn){
	ws.mutex.Lock()
	delete(ws.Clients,conn)
	ws.mutex.Unlock()
}

func (ws *Websocket) PushToFilteredClients(filter func(*websocket.Conn) bool, msg string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	for client := range ws.Clients {
		if !filter(client) {
			continue
		}

		err := client.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			switch e := err.(type) {
			case *websocket.CloseError:
				log.Printf("Client closed connection: %v (code: %d)", e.Text, e.Code)
			case net.Error:
				if e.Timeout() {
					log.Println("Timeout writing to client:", e)
				}else {
					log.Println("Network error:", e)
				}
			default:
				log.Printf("Unexpected WebSocket error: %v", err)
			}

			client.Close()
			delete(ws.Clients, client)
		}
	}
}
