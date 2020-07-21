package main

import (
	"log"

	"github.com/antihax/AtlasNotificationViewer/mapserver"
)

func main() {
	s := mapserver.NewMapServer()
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Server quit!")
}
