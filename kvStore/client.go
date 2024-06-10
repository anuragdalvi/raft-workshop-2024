package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func startClient() {
	// Define server address
	serverAddr := &net.TCPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 12343,
	}

	// Connect to the TCP server
	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("KV>> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		response, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		fmt.Print("Server response: " + response)

	}

}
