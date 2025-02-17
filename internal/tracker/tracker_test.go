package tracker

import (
	"bytes"
	"net"
	"reflect"
	"testing"
	"time"
)

type mockConn struct {
	input  *bytes.Buffer
	output *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (n int, err error)   { return m.input.Read(b) }
func (m *mockConn) Write(b []byte) (n int, err error)  { return m.output.Write(b) }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (m *mockConn) RemoteAddr() net.Addr               { return &net.TCPAddr{} }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestStartServer(t *testing.T) {
	tracker := &Tracker{
		Address: "localhost:6000",
		Peers:   make(map[string]map[string]string),
	}
	go tracker.StartServer()
	time.Sleep(100 * time.Millisecond) // Give server time to start

	conn, err := net.Dial("tcp", "localhost:6000")
	if err != nil {
		t.Fatalf("Failed to connect to tracker server: %v", err)
	}
	conn.Close()
}

func TestHandleConnection(t *testing.T) {
	tracker := &Tracker{
		Peers: make(map[string]map[string]string),
	}

	registerMessage := "REGISTER peer1:9090\nsong1 hash1\nsong2 hash2\n\n"
	mock := &mockConn{input: bytes.NewBufferString(registerMessage), output: new(bytes.Buffer)}
	tracker.handleConnection(mock)

	if mock.output.String() != "OK\n" {
		t.Errorf("Expected 'OK', got '%s'", mock.output.String())
	}
}

func TestRegister(t *testing.T) {
	tracker := Tracker{
		Address: "localhost:8080",
		Peers:   make(map[string]map[string]string),
	}

	peerAddress := "peer1:9090"
	songs := map[string]string{
		"song1": "hash1",
		"song2": "hash2",
	}

	tracker.Register(peerAddress, songs)

	if !reflect.DeepEqual(tracker.Peers[peerAddress], songs) {
		t.Errorf("Expected %v, got %v", songs, tracker.Peers[peerAddress])
	}
}

func TestQuery(t *testing.T) {
	tracker := Tracker{
		Address: "localhost:8080",
		Peers: map[string]map[string]string{
			"peer1:9090": {"song1": "hash1"},
			"peer2:9091": {"song2": "hash2"},
			"peer3:9092": {"song1": "hash3"},
		},
	}

	expected := []string{"peer1:9090", "peer3:9092"}
	result := tracker.Query("song1")

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestQueryNoResults(t *testing.T) {
	tracker := Tracker{
		Address: "localhost:8080",
		Peers: map[string]map[string]string{
			"peer1:9090": {"song1": "hash1"},
			"peer2:9091": {"song2": "hash2"},
		},
	}

	result := tracker.Query("song3")

	if len(result) != 0 {
		t.Errorf("Expected empty list, got %v", result)
	}
}
