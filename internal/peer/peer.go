package peer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Peer struct {
	Address    string
	Songs      map[string]string
	TrackerCon net.Conn
}

func (p *Peer) RegisterWithTracker() {
	var payload []byte
	payload = append(payload, []byte("REGISTER "+p.Address+"\n")...)
	for name, loc := range p.Songs {
		payload = append(payload, []byte(name+" "+loc+"\n")...)
	}
	payload = append(payload, []byte("\n\n")...)

	_, err := p.TrackerCon.Write(payload)
	if err != nil {
		log.Printf("Error sending REGISTER request: %v", err)
		return
	}

	reader := bufio.NewReader(p.TrackerCon)
	response, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Error reading response from tracker: %v", err)
		return
	}

	if response == "OK\n" {
		log.Println("Successfully registered with the tracker.")
	} else {
		log.Printf("Failed to register with tracker: %s", response)
	}
}

// Listen for requests from clients
func (p *Peer) StartServer() {

	listener, err := net.Listen("tcp", p.Address)

	if err != nil {
		log.Fatalf("Error starting peer server: %v \n", err)
	}
	fmt.Println("Peer server running")
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalf("Error accepting peer server connection: %v \n", err)
		}

		go p.handleConnection(conn)
	}

}

func (p *Peer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')

	if err != nil {
		fmt.Fprintf(conn, "Error reading peer request: %v", err)
		return
	}
	msg = strings.TrimSpace(msg)

	parts := strings.SplitN(msg, " ", 2)

	command := parts[0]

	switch command {
	case "REQUEST_SONG":
		song := parts[1]

		if path, ok := p.Songs[song]; ok {

			if !fileExists(path) {
				delete(p.Songs, song)
				//TODO: send an UPDATE request to tracker to remove the song from tracker
				conn.Write([]byte("Song does not exist"))
				return
			}
			m3u8 := getM3U8Path(path)
			if !fileExists(m3u8) {
				path = convertMP3(path, song)
			}
			_, err = fmt.Fprintf(conn, "PLAYLIST %s\n", song)
			if err != nil {
				fmt.Fprintf(conn, "Failed to write header: %v", err)
			}
			sendFile(conn, path)
		}
	case "REQUEST_SEGMENT":
		segment := parts[1]
		path, err := getSegment(segment)
		if err != nil {
			//TODO: improve this error handling
			log.Fatal("Error finding file")
		}
		_, err = fmt.Fprintf(conn, "SEGMENT %s\n", segment)
		if err != nil {
			fmt.Fprintf(conn, "Failed to write header: %v", err)
		}
		sendFile(conn, path)
	}
}
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func sendFile(conn net.Conn, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	_, err = io.Copy(conn, file)
	if err != nil {
		return fmt.Errorf("failed to send file: %v", err)
	}

	return nil
}

func getM3U8Path(songPath string) string {
	ext := filepath.Ext(songPath)
	return strings.TrimSuffix(songPath, ext) + ".m3u8"
}

func convertMP3(path string, songName string) string {
	wd, _ := os.Getwd()
	playlistDir := fmt.Sprintf("%s/playlist", wd)
	if _, err := os.Stat(playlistDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(playlistDir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	playlistPath := fmt.Sprintf("%s/%s.m3u8", playlistDir, songName)
	fileLoc, err := filepath.Abs(path)

	if err != nil {
		//TODO: improfve this error handling
		log.Fatal("Error finding file")
	}

	fmt.Printf("Playlist Location: %s \n", playlistPath)
	fmt.Printf("Song File Location: %s \n", fileLoc)

	cmd := exec.Command(
		"ffmpeg",
		"-i",
		fileLoc,
		"-c:a",
		"aac",
		"-b:a",
		"128k",
		"-f",
		"hls",
		"-hls_time",
		"6",
		"-hls_playlist_type",
		"event",
		playlistPath,
	)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}

	fmt.Println("Command executed successfully")

	return playlistPath
}

func getSegment(segment string) (string, error) {
	path, err := filepath.Abs("./playlist/" + segment)
	return path, err

}
