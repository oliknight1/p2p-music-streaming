package main

import "p2p-music/internal"

func main() {
	t := internal.Tracker{
		Address: ":8000",
		Peers:   make(map[string]map[string]string),
	}
	t.StartServer()
}
