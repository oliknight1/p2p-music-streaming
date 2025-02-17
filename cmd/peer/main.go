package main

import (
	"log"
	"net"
	peer "p2p-music/internal/peer"
	"path/filepath"
)

func main() {
	songPath, _ := filepath.Abs("./SongA.mp3")
	p := peer.Peer{
		Address: ":8001",
		Songs: map[string]string{
			"SongA": songPath,
		},
		TrackerCon: nil,
	}

	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Fatal("Error connecting to tracker: ", err)
	}

	p.TrackerCon = conn

	//TODO: Return error?
	p.RegisterWithTracker()

	p.StartServer()
}
