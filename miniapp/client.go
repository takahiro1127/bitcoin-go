package miniapp

import (
    "fmt"
    "net"
    "time"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:50030")
    if err != nil {
        fmt.Printf("Dial error: %s\n", err)
        return
    }
    defer conn.Close()

    sendMsg := "Hello! This is test message from my sample client.\n"
    // sendMsg := "stop"
    conn.Write([]byte(sendMsg))
    readBuf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	readlen, err := conn.Read(readBuf)
	fmt.Println("server: " + string(readBuf[:readlen]))
}