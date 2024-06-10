package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type raftNet struct {
	// My Address
	nodeNum int
}

// Send a message to a specific node number (returns immediately)
func (*raftNet) send(destination int, message string) {
	// Get Port
	destinationPort, _ := strconv.Atoi(getConfiguration()[destination][1])

	// Define server address
	serverAddr := &net.TCPAddr{
		IP:   net.ParseIP("localhost"),
		Port: destinationPort,
	}

	// Connect to the TCP server
	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// send Message
	_, err = conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Error writing to server:", err)
		return
	}

}

// Receive and return any message sent to me (blocks)
func (*raftNet) receive(nodeNum string) string {
	node, _ := strconv.Atoi(nodeNum)
	nodePort := getConfiguration()[node][1]
	source := "localhost:" + nodePort
	listener, err := net.Listen("tcp", source)
	if err != nil {
		return "err"
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return "err"
	}

	reader := bufio.NewReader(conn)
	input, err := reader.ReadString('\n')
	message := strings.TrimSuffix(input, "\n")
	return message
}
