package main

import ( "github.com/gorilla/websocket"
		 "github.com/gocraft/web"
		 "errors"
		 "strings"
		 "fmt"
		 "time"
		 "encoding/json")

type Client struct {
	addr string
	ws   *websocket.Conn
}

type WebData struct {
	Datatype string `json:"type"`
	Msg      string `json:"msg"`
	Ip       string `json:"ip"`
}

var clients map[string] *Client

var webChannel chan WebData

var logContextWs LoggerContext
var wsLocalStatus *StatusModule

var upgrader = websocket.Upgrader {
    	ReadBufferSize: 1024,
    	WriteBufferSize: 1024,
}

func NewClient(w web.ResponseWriter, r *web.Request) error {
	addr := strings.Split(r.RemoteAddr, ":") [0]
	conn, err := upgrader.Upgrade(w, r.Request, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("websocket connection failed, %s", err))
	}

	client := &Client{
		addr: addr,
		ws:   conn,
	}
	clients[addr] = client
	return nil
}

func RemoveClient(ip string) error{
	err := clients[ip].ws.Close()
	if err != nil{
		return errors.New(fmt.Sprintf("Error while closing websocket from %s, %s",ip, err))
	}
	delete(clients, ip)
	return nil
}

func Loop(){
	previous := ""
	for {

		previous = CheckServerStatus(previous)

  		select {
  			case m := <- webChannel:
  					Send(m)
  			default: time.Sleep(50 * time.Millisecond)
  		}
	}
}
func CheckServerStatus(previous string) string{
	current, _ := wsLocalStatus.GetState()

	if current != previous {
		msg := WebData{Datatype: "serverstatus", Msg: current}
		js, err :=  json.Marshal(msg)
		if err != nil {
        LogErr(logContextWs, "json creation failed")
    }

    for _, c := range clients {
	    err = c.ws.WriteMessage(websocket.TextMessage, js)
	    if err != nil {
	        LogErr(logContextWs, "websocket message sending failed, %s", err)
	    }
	}
	previous = current
	}
	return previous
}

func Send(m WebData) {
	json, err := json.Marshal(m)
	if err != nil {
		LogErr(logContextWs, "json creation failed")
	}

	for ip, c := range clients {
		if m.Ip==ip {
				err = c.ws.WriteMessage(websocket.TextMessage, json)
				if err != nil {
	    		LogErr(logContextWs, "websocket message sending failed, %s", err)
			}
		}
	}
	time.Sleep(100 * time.Millisecond)
}

//WebsocketInit...
func WebsocketInit(sm *StatusModule) {
	wsLocalStatus = sm
	logContextWs = LoggerContext{
		name:  "WEBSOCKET",
		level: 3}

	webChannel = make(chan WebData)

	clients = make(map[string]*Client)
	go Loop()
}