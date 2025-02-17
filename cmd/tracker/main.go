package main

import tracker "p2p-music/internal/tracker"

func main() {
	t := tracker.Tracker{
		Address: ":8000",
		Peers:   make(map[string]map[string]string),
	}

	t.StartServer()
}
