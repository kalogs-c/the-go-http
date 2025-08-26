package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
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
	file, err := os.Open("./messages.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	contentsCh := getLinesChannel(file)
	for line := range contentsCh {
		fmt.Printf("read: %s\n", line)
	}
}
