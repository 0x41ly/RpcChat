package main

import (
	commons "RpcChat/Commons"
	"flag"
	"log"
	"net"
	"net/rpc"
)

func parseFlags(cs *commons.ChatServer) {
	flag.StringVar(&cs.Port, "Port", "3410", "Port for chat server to listen on")
	flag.Parse()

	cs.Port = ":" + cs.Port
}

func RunServer(cs *commons.ChatServer) {
	rpc.Register(cs)

	log.Printf("Listening on Port %s...\n", cs.Port)

	l, err := net.Listen("tcp", cs.Port)
	if err != nil {
		log.Panicf("Can't bind Port to listen. %q", err)
	}

	rpc.Accept(l)
}

func main() {
	cs := new(commons.ChatServer)
	cs.MessageQueue = make(map[string][]string)
	cs.KeepMeAlive = make(chan bool, 1)

	parseFlags(cs)
	RunServer(cs)

	<-cs.KeepMeAlive
}
