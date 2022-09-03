package controller

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

/*participant decribes a single entity in the hashmap*/
type Participant struct {
	Host bool
	Conn *websocket.Conn
}

/*RoomMap is the main hashmap [roomID string] -> [[]Participant]*/
type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}

/*Init initialize the RoomMap struct*/
func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

/*Get will return the array of participants in the room*/
func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

/*CreateRoom generate a unique ID and return it  -> insert it in the hashmap*/
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	//generate random number or use google -> uuid
	rand.Seed(int64(time.Now().Nanosecond()))
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	roomID := string(b)
	r.Map[roomID] = []Participant{}

	return roomID
}

/* will create a participant and andd it in the hashmap*/
func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{host, conn}

	log.Println("Inserting into into Room with RoomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

/*Delete a room with room id*/
func (r *RoomMap) DeleteRoom(RoomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, RoomID)
}
