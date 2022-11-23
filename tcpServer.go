package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var Map sync.Map

// listening for the established connection
func RunTCPServer(tcp_addr string) {
	l, err := net.Listen("tcp", tcp_addr)
	if err != nil {
		log.Println("Error appeared at listening.", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error appeared at accepting.", err.Error())
			os.Exit(1)
		}
		// we handle connections in a new goroutine
		go handleRequest(conn, tcp_addr)
	}
}

// Due to requests we will store ans get responses
func handleRequest(conn net.Conn, tcp_addr string) {
	d := json.NewDecoder(conn)
	var msg DataRequest
	err := d.Decode(&msg)
	if err != nil {
		log.Fatal("Error appeared at decoding the TCP server json format.")
	}

	// In dependence to what operations we have
	// we will store ans get responses
	var resp string
	switch msg.Operation {

	//In case of CREATE
	case "POST":
		Map.Store(msg.Key, msg.Val)
		resp = fmt.Sprintf("Created at %v", tcp_addr)

		//In case of UPDATE
	case "PUT":
		_, ok := Map.Load(msg.Key)
		if ok {
			Map.Store(msg.Key, msg.Val)
			resp = fmt.Sprintf("Updated at %v", tcp_addr)
		}

		//In case of READ
	case "GET":
		storeVal, ok := Map.Load(msg.Key)
		if ok {
			resp = fmt.Sprint(storeVal)
		} else {
			resp = "There is nothing like this."
		}

		//In case of DELETE
	case "DELETE":
		_, ok := Map.Load(msg.Key)
		if ok {
			Map.Delete(msg.Key)
			resp = fmt.Sprintf("Deleted at %v", tcp_addr)
		}

	}
	conn.Write([]byte(resp))
	conn.Close()
}
