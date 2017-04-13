package main

import (
	"log"
	"net/http"
	"os"
)

type msg struct {
	Num int
}

type logFuncs interface {
	Println(v ...interface{})
	Fatal(v ...interface{})
	Printf(format string, v ...interface{})
}

var logging logFuncs = log.New(os.Stdout, "[ ] ", log.LstdFlags)

func main() {
	hub := NewHub()
	go hub.Run()

	r := http.NewServeMux()
	r.Handle("/", http.FileServer(http.Dir("./react-chat-client/build/"))) //handles static html / css etc
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { serveWs(hub, w, r) })
	http.Handle("/", r)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logging.Fatal("Http server fell down with: ", err)
	}

}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Println(err)
		return
	}
	client := NewClient(hub, conn)
	client.hub.Register <- client
	go client.writePump()
	client.readPump()
}
