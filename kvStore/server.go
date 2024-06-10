package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func run_command(command string, kv Mapping) string {

	cmd := strings.Split(command, " ")

	cmdInit := strings.ToLower(cmd[0])

	switch cmdInit {
	case "get":
		return kv.Get(cmd[1])
	case "set":
		return kv.Set(cmd[1], cmd[2])
	case "delete":
		return kv.Delete(cmd[1])
	case "snapshot":
		return kv.Snapshot(cmd[1])
	case "retrieve":
		return kv.RetrieveSnapshot(cmd[1])

	default:
		return "Invalid Command"
	}

}

func startServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	kv := NewKeyValeStore()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go handleConnection(conn, *kv)
	}
}

func handleConnection(conn net.Conn, kv Mapping) {

	reader := bufio.NewReader(conn)
	for {
		input, err := reader.ReadString('\n')
		fmt.Printf("Input: %s", input)
		message := strings.TrimSuffix(input, "\n")
		if err == io.EOF {
			// Client disconnected gracefully
			conn.Close()
		}
		if message == "exit" {
			conn.Close()
		}

		commandOutput := run_command(message, kv)
		cOutput := commandOutput + "\n"
		fmt.Printf("Output: %s", cOutput)
		conn.Write([]byte(cOutput))
	}

}
