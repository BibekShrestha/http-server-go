package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	cwd, _ := os.Getwd()
	var directory = flag.String("directory", cwd, "Working directory for the http server")
	flag.Parse()

	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	cxt := context.WithValue(context.Background(), "workDir", *directory)
	fmt.Println("Working from workdir: ", *directory)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(cxt, conn)
	}
}
