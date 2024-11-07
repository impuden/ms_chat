package main

import (
	"chat-service/config"
	"chat-service/db"
	"chat-service/handlers"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var frontaddr = flag.String("addr", ":1200", "http service address")
var coreAddr = flag.String("addrback", ":1201", "http service address")

func main() {
	config.LoadConfig()

	logFile, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	db.DB, err = db.ConnectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД ", err)
	}
	defer db.DB.Close()

	//flag.Parse()
	hub := handlers.NewHub()
	defer func() {
		for roomName := range hub.Rooms {
			hub.RemoveRoom(roomName)
		}
	}()

	go func() {
		http.HandleFunc("/apicore", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerCore(w, r)
		})

		log.Println("server for backend is ready", *coreAddr)
		server := &http.Server{
			Addr:              *coreAddr,
			ReadHeaderTimeout: 3 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
			log.Fatal("backend server failed: ", err)
		}
	}()

	go func() {
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleFront(hub, w, r)
		})

		log.Println("front server was start", *frontaddr)
		server := &http.Server{
			Addr:              *frontaddr,
			ReadHeaderTimeout: 3 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
			log.Fatal("fail: ", err)
		}
		log.Println("front server stop")
	}()

	select {}
}
