package main

import (
    "net"
    "log"
    "os"
    "io"
    "sync"
)

var wg sync.WaitGroup

const (
    PORT = ":1052"
)

func Forward(conn net.Conn, toAddr string) {
    wg.Add(1)
    client, err := net.Dial("tcp", toAddr)
    connRemoteAddr := conn.RemoteAddr().String()
    if err != nil {
        log.Printf("type=error, value=Unable to open connection with the target TCP server\n")
    }
    log.Printf("type=success, value=Connected to target server %s\n", toAddr)
    go func() {
        defer client.Close()
        defer conn.Close()
        io.Copy(client, conn)
    }()
    go func() {
        defer client.Close()
        defer conn.Close()
        io.Copy(conn, client)
    }()
    log.Printf("type=success, value=Successfully forwarded connection from %s to %s", connRemoteAddr, toAddr)
    wg.Done()
}

func main() {
    if len(os.Args) == 2 {
        forwardTo := os.Args[1]
        listener, err := net.Listen("tcp", PORT)
        if err != nil {
            log.Printf("type=error, value=Unable to listen for connections on port %s\n", PORT)
        }
        log.Println("type=success, value=Listening for connections on port %s", PORT)
        for {
            conn, err := listener.Accept()
            if err != nil {
                log.Printf("type=error, value=Unable to accept incoming connection from %s\n", conn.RemoteAddr().String())
            }
            go Forward(conn, forwardTo)
        }
    } else {
        log.Println("type=error, value=Incorrect Usage")
        log.Println("\tCorrect Usage: <server ip:port>")
    }
    wg.Wait()
}
