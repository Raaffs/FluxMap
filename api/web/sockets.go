package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSUser struct{
	Username 	string
	Connected 	bool
}

type Websocket struct{
	Upgrader websocket.Upgrader
	//conn=>username 
	Clients  map[*websocket.Conn]WSUser
	mutex 	 sync.Mutex
	Msg 	 chan []byte
}


func NewWS()*Websocket{
	ws:=&Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool{
				return true
			},
		},
		Clients: make(map[*websocket.Conn]WSUser),
		mutex: sync.Mutex{},
	}
	return ws
}

func (ws *Websocket)AddClient(conn *websocket.Conn, username string){
	ws.mutex.Lock()
	ws.Clients[conn]=WSUser{
		Username: username,
		Connected: true,
	}
	ws.mutex.Unlock()
}

func (ws *Websocket)RemoveClient(conn *websocket.Conn){
	ws.mutex.Lock()
	delete(ws.Clients,conn)
	ws.mutex.Unlock()
}

func (ws *Websocket) CloseAllClients() {
    ws.mutex.Lock()
    defer ws.mutex.Unlock()

    for client := range ws.Clients {
        client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutdown"))
        client.Close()
        delete(ws.Clients, client)
    }
}

func (ws *Websocket)KeepAlive(w http.ResponseWriter, r *http.Request, username string){
	conn,err:=ws.Upgrader.Upgrade(w,r,nil);if err!=nil{
		log.Println("Error connecting upgrading http connection: ",err)
		return
	}

	ws.AddClient(conn, username)

	defer func(){
		ws.RemoveClient(conn)
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(appData string)error{
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
	
	ticker:=time.NewTicker(30*time.Second)
	defer ticker.Stop()

	for range ticker.C{
		if err:=conn.WriteMessage(websocket.PingMessage,[]byte{}); err!=nil{
				log.Println("Error writing message to client: ",err)
				break
		}
	}
}

// PushToClients broadcasts a message to matching clients.
// It returns a channel of errors from failed client sends.
// IMPORTANT: The caller **must** read from the returned channel to avoid goroutine leaks!
func (ws *Websocket)PushToClients(filter func(WSUser) bool, msg []byte) <-chan error {
	ws.mutex.Lock()
	var wg sync.WaitGroup
	var clients []*websocket.Conn
	for clientConn,user:= range ws.Clients{
		if filter(user){
			clients=append(clients, clientConn)
		}
	}
	errchan:=make(chan error, len(clients))
	ws.mutex.Unlock()
	
	for _,client:=range clients{
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err:=client.WriteMessage(websocket.TextMessage,msg); err!=nil{
				log.Println("Error writing message to client: ",err)
				errchan<-err
				ws.mutex.Lock()
				delete(ws.Clients,client)
				ws.mutex.Unlock()
			}
		} ()
	}
	go func ()  {
		wg.Wait()
		close(errchan)	
	}()
	return errchan
}

func (ws *Websocket) detectClientDisconnect(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error or client disconnected:", err)
			ws.RemoveClient(conn)
			conn.Close()
			break
		}
	}
}