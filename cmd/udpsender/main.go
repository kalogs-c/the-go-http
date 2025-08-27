package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		message, err := reader.ReadString(byte('\n'))
		if err != nil {
			panic(err)
		}

		if _, err = conn.Write([]byte(message)); err != nil {
			log.Println(err)
		}
	}
}
