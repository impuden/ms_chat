package handlers

import (
	"chat-service/auth"
	"chat-service/db"
	"chat-service/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandlerCore(w http.ResponseWriter, r *http.Request) {
	staticToken := os.Getenv("APP_TOKEN")
	reqToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(reqToken, "Bearer ") || reqToken[7:] != staticToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var coreData models.CoreData
	if err := json.NewDecoder(r.Body).Decode(&coreData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fmt.Println(coreData)

	roomID, err := db.CheckCreateRoom(db.DB, coreData.User1, coreData.User2, coreData.ItemID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jwtToken, err := auth.GenerateToken(roomID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Contant-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}

func HandleFront(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("tryin")

	jwtToken := r.Header.Get("Authorization")
	fmt.Println("token come: ", jwtToken)
	if jwtToken == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		fmt.Println("error with jwt: ")
		return
	}
	log.Println("jwt received: ", jwtToken)

	RoomID, err := auth.ParseToken(jwtToken)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		log.Println("error unauth token: ", err)
		return
	}
	log.Println("pars is ok: ", RoomID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	var room *Room
	var client *Client
	clientRegistered := false
	offset := 0
	const initialLoad = 40
	const limit = 15

	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("error read json: ", err)
			break
		}
		log.Println("msg recieved: ", msg)

		if err := db.SaveMessage(db.DB, msg.Username, RoomID, msg.Message); err != nil {
			log.Println("error saving message: ", err)
		}

		if !clientRegistered {
			room = hub.GetRoom(RoomID)
			if room == nil {
				log.Println("room not found: ", err)
				break
			}

			client = &Client{Room: room, Conn: conn, Send: make(chan []byte, 256)}
			client.Room.Register <- client

			go client.WritePump()
			go client.ReadPump()

			messages, err := db.LoadMessages(db.DB, RoomID, 0, 40)
			if err != nil {
				log.Println("error load messages: ", err)
			} else {
				for _, msg := range messages {
					err := conn.WriteJSON(msg)
					if err != nil {
						log.Println("error write json message: ", err)
						continue
					}

				}
				log.Println("1")
			}
			clientRegistered = true
			offset = initialLoad
		} else if msg.Message == "LOAD_MORE" {
			messages, err := db.LoadMessages(db.DB, RoomID, offset, limit)
			if err != nil {
				log.Println("err load msg: ", err)
			} else {
				for _, msg := range messages {
					jsonMsg, err := json.Marshal(msg)
					if err != nil {
						log.Println("err marsh msg: ", err)
						continue
					}
					client.Send <- jsonMsg
				}
			}
			offset += limit
		} else {
			log.Println("2")
			if err := db.SaveMessage(db.DB, msg.Username, RoomID, msg.Message); err != nil {
				log.Println("error saving message: ", err)
			}

			broadcastMessage := models.Message{Username: msg.Username, Message: msg.Message, Timestamp: msg.Timestamp}
			bcMsg, err := json.Marshal(broadcastMessage)
			if err != nil {
				log.Println("error marshal broadcast message: ", err)
				break
			}
			room.Broadcast <- bcMsg
		}
	}
}
