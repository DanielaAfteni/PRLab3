package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type DataRequest struct {
	Operation string `json:"operation"`
	Val       string `json:"val"`
	Key       string `json:"key"`
}

var mapOfServers map[int]string

func GetRouter(m map[int]string) *mux.Router {
	mapOfServers = m
	r := mux.NewRouter()

	r.HandleFunc("/create/{key}/{value}", PostValue).Methods("POST")
	r.HandleFunc("/read/{key}", GetValue).Methods("GET")
	r.HandleFunc("/update/{key}/{value}", UpdateValue).Methods("PUT")
	r.HandleFunc("/delete/{key}", DeleteValue).Methods("DELETE")

	return r
}

func DialTCPServer(tcp_addr string, key, val, operation string) string {
	conn, err := net.Dial("tcp", tcp_addr)
	if err != nil {
		fmt.Println("Error appeared.", err)
	}
	messageStructure := DataRequest{Operation: operation, Key: key, Val: val}
	msg, _ := json.Marshal(messageStructure)
	fmt.Fprint(conn, string(msg))
	message, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	return message

}

func GetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	var resp string
	for i := 1; i <= len(mapOfServers); i++ {
		resp = DialTCPServer(mapOfServers[i], key, "NONE", "GET")
		if resp != "There is nothing like this." {
			break
		}
	}
	fmt.Fprint(w, resp)
}

func DeleteValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	var resp string
	for i := 1; i <= len(mapOfServers); i++ {
		resp += DialTCPServer(mapOfServers[i], key, "NONE", "DELETE")
	}
	fmt.Fprint(w, resp)
}

func UpdateValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["value"]
	var resp string
	for i := 1; i <= len(mapOfServers); i++ {
		resp = DialTCPServer(mapOfServers[i], key, val, "PUT")
	}
	fmt.Fprint(w, resp)
}

type RoundBoutCounter struct {
	m sync.Mutex
	c int
}

var rbc = RoundBoutCounter{sync.Mutex{}, 1}

func PostValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["value"]
	rbc.m.Lock()
	if rbc.c%(len(mapOfServers)+1) == 0 {
		rbc.c = 1
	}
	temp := rbc.c
	rbc.c += 1
	rbc.m.Unlock()
	var resp string
	for i := 0; i < int(len(mapOfServers)/2+1); i++ {
		resp += DialTCPServer(mapOfServers[temp], key, val, "POST")
		temp += 1
		if temp%(len(mapOfServers)+1) == 0 {
			temp = 1
		}
	}
	fmt.Fprint(w, resp)
}
