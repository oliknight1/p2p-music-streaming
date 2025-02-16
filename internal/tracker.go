package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Tracker struct {
	Address string
	Peers   map[string]map[string]string
	mu      sync.Mutex
}

func (t *Tracker) Register(peerAddress string, songs map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Peers[peerAddress] = songs
}

// Return a list of peer addresses that have the requested song
// Will return an empty list if no peers have the song
func (t *Tracker) Query(songName string) []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	var peers []string
	for peer, songs := range t.Peers {
		if _, ok := songs[songName]; ok {
			peers = append(peers, peer)
		}
	}
	return peers
}

func (t *Tracker) StartServer() {
	listener, err := net.Listen("tcp", t.Address)

	if err != nil {
		log.Fatalf("Error starting tracker server: %v \n", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalf("Error accepting tracker server connection: %v \n", err)
		}

		go t.handleConnection(conn)
	}

}

func (t *Tracker) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')

	if err != nil {
		fmt.Fprintf(conn, "Error reading tracker request: %v ", err)
	}

	msg = strings.TrimSpace(msg)

	parts := strings.SplitN(msg, " ", 2)

	command := parts[0]

	switch command {

	case "REGISTER":
		address := parts[1]
		songs := make(map[string]string)
		for {
			line, err := reader.ReadString('\n')

			if len(parts) < 2 {
				conn.Write([]byte("ERROR Invalid REGISTER request\n"))
			}
			if err != nil {
				fmt.Fprintf(conn, "Error reading tracker request: %v ", err)
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}

			parts := strings.Split(line, " ")
			songs[parts[0]] = parts[1]
		}
		t.Register(address, songs)
		conn.Write([]byte("OK\n"))

	case "QUERY":
		if len(parts) < 2 {
			conn.Write([]byte("ERROR Invalid QUERY request\n"))
			return
		}
		songName := parts[1]
		peers := t.Query(songName)

		if len(peers) == 0 {
			conn.Write([]byte("PEERS \n"))
		}
		response := "PEERS " + strings.Join(peers, " ") + "\n"

		conn.Write([]byte(response))

	default:
		log.Fatal("Unknown command")
	}
}
