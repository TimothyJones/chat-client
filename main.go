package main

import (
	"fmt"
	"io/ioutil"
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

var logging logFuncs

func main() {
	hub := NewHub()
	go hub.Run()

	r := http.NewServeMux()
	//	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./server/public/")))
	r.Handle("/", http.FileServer(http.Dir("./react-chat-client/build/"))) //handles static html / css etc. under ./webroot
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { serveWs(hub, w, r) })
	http.Handle("/", r)

	logging = log.New(os.Stdout, "[ ]", log.LstdFlags)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logging.Fatal("Http server fell down with: ", err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
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
