package websocket_gateway

import (
	"application/libraries/utils"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  5120,
	WriteBufferSize: 5120,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 5120
)

type ClientInfo struct {
	Verson   string
	Platform string
}

type Client struct {
	hub *GatewayInfo

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	Id string

	Num int

	Info ClientInfo

	Idfa string
}

// Client is a middleman between the websocket connection and the hub.

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {

	defer func() {
		//c.hub.unregister <- c
		utils.Logger().Println("ReadPump error")
		c.hub.AddUnregister(c)
		if err := recover(); err != nil {
			utils.Logger().Printf("panic: %v\n", err)
		}

	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mType, message, err := c.conn.ReadMessage()
		fmt.Println("99",message,err)

		if err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				utils.Logger().Println("error: websocket is close client idfa :", c.Idfa, "error message:", err.(*websocket.CloseError).Error(), "mtype:", mType)
				return
			}

			utils.Logger().Println("error: ", err, "idfa:", c.Idfa, "mtype:", mType, "message:", string(message))
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		utils.Logger().Printf("message is : %s \n", message)
		GetRuteInstance().Serve(c, message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {

	defer func() {
		//c.hub.unregister <- c
		utils.Logger().Println("WritePump error")
		c.hub.AddUnregister(c)
		if err := recover(); err != nil {
			utils.Logger().Printf("panic: %v\n", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				utils.Logger().Println("NextWriter error", err, " idfa:", c.Idfa)
				return
			}

			fmt.Printf("WritePump:%s\n", message)
			w.Write(message)
			/*
				// Add queued chat messages to the current websocket message.
				n := len(c.send)
				fmt.Printf("n%\n", n)
				for i := 0; i < n; i++ {
					w.Write(newline)
					w.Write(<-c.send)
				}
			*/
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) PingPump() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		//c.hub.unregister <- c
		c.hub.AddUnregister(c)
		if err := recover(); err != nil {
			utils.Logger().Printf("panic: %v\n", err)
		}
	}()
	for {
		<-ticker.C

		if _, ok := c.hub.ClientIdToUid.Load(c.Id); !ok {
			return
		}
		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			c.Num++
		}

		c.Num = 0

		if c.Num >= 3 {
			return
		}
	}
}

func (c *Client) Replay(message []byte) {
	c.send <- message
	duration := time.Millisecond * 20
	time.Sleep(duration)

}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *GatewayInfo, w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	utils.Logger().Println("connect IP is ", conn.RemoteAddr().String())
	connId := fmt.Sprintf("%d:%s", time.Now().Nanosecond, conn.RemoteAddr().String())

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), Id: utils.Md5(connId)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
	go client.PingPump()
}
