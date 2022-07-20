package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/LUXROBO/server/libs/env"
	"github.com/LUXROBO/server/libs/logger"
	"github.com/LUXROBO/server/services/chat/internal/service"
)

const serviceName = "chat"

func main() {
	env := env.Initialize(serviceName)
	logger.Initialize(env.ServiceName, env.Mode)
	hub := initHub()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/health", health)
	http.HandleFunc(env.ChatWebsocketPath, func(w http.ResponseWriter, r *http.Request) {
		service.ServerWs(hub, w, r)
	})

	addr := flag.String("http", ":"+env.Port, "HTTP listen address")
	logger.Infof("connect to http://localhost%s for %s server", *addr, serviceName)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initHub() *service.Hub {
	hub := service.NewHub()
	go hub.Run()
	go hub.Subscribe()
	return hub
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	logger.Debug(r.URL)

	if r.URL.Path == "/test" {
		http.ServeFile(w, r, "./demo/index2.html")
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "./demo/index.html")
}

func health(w http.ResponseWriter, r *http.Request) {
	response := "good"
	_, error := w.Write([]byte(response))
	if error != nil {
		panic(error)
	}
}
