package Revenge

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
)

type wsocket_connect interface{}

type wsocket struct {
	sock_id uint32
}

func (ws *wsocket) start() {
	fmt.Println("From websocket initializer")
}

func generateWebSocketAcceptKey(clientKey string) string {
	const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(clientKey + magicGUID))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, _ := http.ReadRequest(reader)

	// Check for a valid WebSocket Upgrade request
	if request.Header.Get("Upgrade") != "websocket" || request.Header.Get("Connection") != "Upgrade" {
		http.Error(conn, "Unsupported protocol", http.StatusBadRequest)
		return
	}

	// Send handshake response
	acceptKey := generateWebSocketAcceptKey(request.Header.Get("Sec-WebSocket-Key"))
	fmt.Fprintf(conn, "HTTP/1.1 101 Switching Protocols\r\n")
	fmt.Fprintf(conn, "Upgrade: websocket\r\n")
	fmt.Fprintf(conn, "Connection: Upgrade\r\n")
	fmt.Fprintf(conn, "Sec-WebSocket-Accept: %s\r\n", acceptKey)
	fmt.Fprintf(conn, "\r\n")

	// Simple echo functionality (no frame handling)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return // end connection if read fails
		}
		fmt.Fprintf(conn, message)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
