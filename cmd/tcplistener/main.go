package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesFromReader(file io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		currentLine := ""
		for {
			content := make([]byte, 8)
			_, err := file.Read(content)
			if err != nil && err != io.EOF {
				panic(err)
			}

			if err == io.EOF {
				close(ch)
				break
			}

			if readAt := bytes.IndexByte(content, '\n'); readAt != -1 {
				currentLine += string(content[:readAt])
				ch <- currentLine

				content = content[readAt+1:]
				currentLine = ""
			}

			currentLine += string(content)
		}
	}()

	return ch
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("connection accepted!")

		for line := range getLinesFromReader(conn) {
			fmt.Println(line)
		}

		fmt.Println("connection closed!")
		conn.Close()
	}
}
