package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"go.bug.st/serial"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal(`Usage: ./main <file> <port>
<file> - For windows use "COMx" for linux use "/dev/ttyx. Globs are supported, but first file will always be used.
<port> - The port to listen on`)
	}

	fs, err := filepath.Glob(args[0])
	if err != nil {
		log.Fatal(err)
	}
	if len(fs) == 0 {
		log.Fatal("No files found")
	}

	mode := &serial.Mode{}
	f, err := serial.Open(fs[0], mode)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	l, err := net.Listen("tcp", ":"+args[1])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer l.Close()

	conns := make(map[net.Conn]bool)
	go ReadFileAndSendToAll(f, conns)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		conns[conn] = true
		go HandleConnection(conn, f, conns)
	}
}

func HandleConnection(c net.Conn, f serial.Port, conns map[net.Conn]bool) {
	defer c.Close()
	b := make([]byte, 1024)
	for {
		n, err := c.Read(b)
		if err != nil {
			break
		}
		f.Write(b[:n])
	}
	delete(conns, c)
}

func ReadFileAndSendToAll(f serial.Port, conns map[net.Conn]bool) {
	fmt.Println("Loading file")
	b := make([]byte, 1024)
	c := 0
	for {
		n, err := f.Read(b)
		if err != nil {
			continue
		}
		fmt.Printf("\r%s\tNum Clients: %d\tNum Mavlink Packets: %d", time.Now().Format("2006-01-02 15:04:05"), len(conns), c)
		c++
		for conn := range conns {
			conn.Write(b[:n])
		}
	}
}
