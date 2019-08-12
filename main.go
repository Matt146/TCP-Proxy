/**********************************************************************************
* By: Matt146
* Date: 2019-08-11 (yyyy/mm/dd)
* Purpose: A TCP Proxy Server
* License: GNU GPL Version 3

Matt146 - TCP Proxy Server
Copyright (C) 2019  Matt146
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
*********************************************************************************/

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

func PrintWelcomeMessage() {
	fmt.Println(`
████████╗ ██████╗██████╗     ██████╗ ██████╗  ██████╗ ██╗  ██╗██╗   ██╗
╚══██╔══╝██╔════╝██╔══██╗    ██╔══██╗██╔══██╗██╔═══██╗╚██╗██╔╝╚██╗ ██╔╝
   ██║   ██║     ██████╔╝    ██████╔╝██████╔╝██║   ██║ ╚███╔╝  ╚████╔╝ 
   ██║   ██║     ██╔═══╝     ██╔═══╝ ██╔══██╗██║   ██║ ██╔██╗   ╚██╔╝  
   ██║   ╚██████╗██║         ██║     ██║  ██║╚██████╔╝██╔╝ ██╗   ██║   
   ╚═╝    ╚═════╝╚═╝         ╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝ 
`)
	fmt.Println("FAQ:")
	fmt.Println("Q1. Does it use SSL?")
	fmt.Println("\tA1: No, the current version does not")
	fmt.Println("How is it licensed?")
	fmt.Println("\tA2: GNU GPL V. 3")
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
