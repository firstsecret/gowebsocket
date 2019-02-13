package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"websocketProject/impl"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("hello"))
	var (
		err    error
		data   []byte
		wsConn *websocket.Conn
		conn   *impl.Connection
	)

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 协程
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}

	}()

	// websocket Conn
	for {
		// Text, Binary
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}

	}
ERR:
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe("0.0.0.0:7777", nil)
}
