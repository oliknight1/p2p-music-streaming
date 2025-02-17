# P2P Music Streaming

## Overview
This project is a peer-to-peer (P2P) music streaming system that allows multiple peers to share and stream music files. A tracker node helps coordinate peers by maintaining a registry of available songs and their respective hosts.

## Features
- **Peer-to-Peer Song Sharing**: Peers register available songs with the tracker and can request songs from other peers.
- **Tracker-Based Discovery**: The tracker node helps peers find the correct source for a song.
- **Chunk-Based Streaming**: Uses `.m3u8` and `.ts` files to facilitate streaming in chunks.
- **Custom Protocol**: Communication between peers and the tracker follows a lightweight text-based protocol.

## Project Structure
```
.
├── cmd
│   ├── tracker     # Entry point for the tracker
│   │   └── main.go
│   ├── peer        # Entry point for peers
│   │   └── main.go
├── internal
│   ├── tracker.go  # Tracker implementation
│   ├── peer.go     # Peer implementation
├── README.md       # Project documentation
├── go.mod          # Go module file
```

## Installation
### Prerequisites
- Go (1.22 or later)
- FFmpeg (for converting audio files into `.m3u8` and `.ts` format)

### Clone the Repository
```sh
git clone https://github.com/yourusername/p2p-music.git
cd p2p-music
```

## Running the Tracker
```sh
cd cmd/tracker
go run main.go
```

## Running a Peer
```sh
cd cmd/peer
```

### **Note**
To run this project locally, you must source an MP3 file for each peer. Ensure that the peer’s song list includes a valid file path to an existing `.mp3` file that will be converted into `.m3u8` and `.ts` formats.

## Custom Protocol
### Registering Songs (Peer → Tracker)
```
REGISTER <peer_address>
SongA /path/to/songA.ts
SongB /path/to/songB.ts

```
**Response:**
```
OK
```

### Querying Songs (Client → Tracker)
```
QUERY SongA
```
**Response:**
```
PEERS peer1_address peer2_address
```

### Requesting a Song (Client → Peer)
```
REQUEST_SONG SongA
```
**Response:**
```
PLAYLIST /path/to/playlist.m3u8
[PLAYLIST_FILE]
```
**Additional request from Client**
```
REQUEST_SEGMENT segment0.ts
```
**Reponse**
```
SEGMENT segment0.ts <binary data>
```

## Architecture Diagram
Below is a high-level overview of how the tracker and peers interact:

```
         +----------------+
         |    Tracker     |
         +----------------+
                ▲
       REGISTER | QUERY
                ▼
+----------------+       +----------------+
|     Peer 1     |       |     Peer 2     |
+----------------+       +----------------+
       ▲  | REQUEST              ▲  |
       |  ▼                      |  ▼
+----------------+       +----------------+
|     Client     |       |     Client     |
+----------------+       +----------------+
```

## Future Enhancements
- Add client playback via `FFplay`

## License
This project is licensed under the MIT License.

