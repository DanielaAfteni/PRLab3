package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Config struct {
	HttpPort            string         `json:"http_port"`
	HttpAddr            string         `json:"http_addr"`
	TcpAddr             string         `json:"tcp_address"`
	PartitionLeaderBool bool           `json:"partition_leader_bool"`
	MapOfServers        map[int]string `json:"tcp_cluster_servers"`
}

func main() {
	conf := GetConfig()
	if conf.PartitionLeaderBool {
		r := GetRouter(conf.MapOfServers)
		go RunTCPServer(conf.TcpAddr)
		log.Printf("Partition leader has the port %v and address %v ", conf.HttpPort, conf.HttpAddr)
		http.ListenAndServe(":"+conf.HttpPort, r)
	} else {
		RunTCPServer(conf.TcpAddr)
	}
}

func GetConfig() *Config {
	jsonFile, err := os.Open("config-1/config.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var c Config
	json.Unmarshal(byteValue, &c)
	return &c
}
