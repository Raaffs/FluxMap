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
	Send	 	chan MessageJob
	Manager		*ConnectionManager
	Websocket	*websocket.Conn
}

type MessageJob struct {
	Message []byte
	Errchan chan error
}

type ConnectionManager struct{
	Upgrader websocket.Upgrader
	Clients  	map[*websocket.Conn]*WSUser
	mutex 	 	sync.Mutex
}

func NewWSUser(username string, conn *websocket.Conn, manager *ConnectionManager) *WSUser {
	return &WSUser{
		Username:  username,
		Connected: true,
		Send:      make(chan MessageJob, 10), // buffered = better for slow clients
		Websocket: conn,
		Manager:   manager,
	}
}

func NewConnectionManager()*ConnectionManager{
	cm:=&ConnectionManager{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool{
				return true
			},
		},
		Clients: make(map[*websocket.Conn]*WSUser),
		mutex: sync.Mutex{},
	}
	return cm
}

func (cm *ConnectionManager)AddClient(user *WSUser){
	cm.mutex.Lock()
	cm.Clients[user.Websocket]=user
	cm.mutex.Unlock()
}

func (cm *ConnectionManager)RemoveClient(conn *websocket.Conn){
	cm.mutex.Lock()
	delete(cm.Clients,conn)
	cm.mutex.Unlock()
}

func (cm *ConnectionManager) CloseAllClients() {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    for client := range cm.Clients {
        client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutdown"))
        client.Close()
        delete(cm.Clients, client)
    }
}



func (client *WSUser)KeepAlive(userWhoSendReq string){
	defer func(){
		client.Manager.RemoveClient(client.Websocket)
		client.Websocket.Close()
	}()

	client.Websocket.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Websocket.SetPongHandler(func(appData string)error{
		return client.Websocket.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	ticker:=time.NewTicker(30*time.Second)
	defer func(){
		ticker.Stop()
		client.Websocket.Close()
	}()

	for{
		select{
		case job,ok:=<-client.Send:
			client.WriteMessage(job,ok)
		case <-ticker.C:
			if err := client.Websocket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

func(client *WSUser)WriteMessage(job MessageJob,ok bool){
	if !ok{
		if err := client.Websocket.WriteMessage(websocket.CloseMessage, nil); err != nil {
			log.Println("connection closed: ", err)
			job.Errchan<-err
		}
		// Return to close the goroutine
		return
	}
	
	if err:=client.Websocket.WriteMessage(websocket.TextMessage,job.Message);err!=nil{
		job.Errchan<-err
		return
	}
}

// PushToClients broadcasts a message to matching clients.
// It returns a channel of errors from failed client sends.
// IMPORTANT: The caller **must** read from the returned channel to avoid goroutine leaks
func (cm *ConnectionManager) PushToClients(msg []byte, filter func(WSUser) bool)<-chan error{
	errchan:=make(chan error)
	go func() {
		cm.mutex.Lock()
		defer cm.mutex.Unlock()
		for _, client := range cm.	Clients {
			if filter(*client){
				client.Send <- MessageJob{
					Message: msg,
					Errchan: errchan,
				}
			}
		}
	}()
	return errchan
}