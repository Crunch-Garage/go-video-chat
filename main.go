package main

import (
	controller "Crunch-Garage/go-video-chat/controllers"
	"log"
	"net/http"
)

func main() {
	controller.AllRooms.Init()

	handleRoutes()
}

func handleRoutes() {

	http.HandleFunc("/create", controller.CreateRoom)
	http.HandleFunc("/join", controller.JoinRoom)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
