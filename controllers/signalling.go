package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*AllRoom is the global hashmap for the controller*/
var AllRooms RoomMap

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type broadcastMessage struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMessage)

func broadcaster() {
	/*whenever a message comes into the broadast channel the roomID is checked and all the
	clients in the particular roomID are fetched
	*/
	for {
		msg := <-broadcast

		for _, client := range AllRooms.Map[msg.RoomID] {
			/*if the client is not equal to the cliet ending the message then send message to them*/
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Fatal(err)
					client.Conn.Close()
				}
			}
		}
	}
}

/*Create a room and return a roomID*/
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)

	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

/*Join the client in a particular room*/
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("roomID is not passed")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("Websocket Upgrade Error :", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	/*create a go routine that will keep listening to a particular channel*/
	go broadcaster()

	for {
		var msg broadcastMessage

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error: ", err)
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		log.Println(msg.Message)

		broadcast <- msg
	}

}
