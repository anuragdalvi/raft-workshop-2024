package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func console(net string) {
	net_number, _ := strconv.Atoi(net)
	network := raftNet{net_number}
	go func() {
		msg := network.receive(net)
		fmt.Println(msg)
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Node %d >> ", net_number)
		input, _ := reader.ReadString('\n')
		destination := strings.Split(input, " ")
		destination_host, _ := strconv.Atoi(destination[0])
		destination_message := strings.TrimSuffix(destination[1], "\n")
		network.send(destination_host, destination_message)
	}

}
