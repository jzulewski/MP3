package server

import (
	"MP3/utils"
	"encoding/gob"
	"net"
)

func handleServer(config utils.Config, ln net.Listener) {
	ch := make(chan utils.Message)
	go handleConnections(ch, ln)
	processIncomingValues(ch, config)
}

func handleConnections(ch chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go unicast_receive(ch, conn)
	}
}

func unicast_receive(ch chan utils.Message, conn net.Conn) {
	for {
		var message utils.Message
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&message)
		utils.CheckError(err)
		ch <- message
	}
}

func processIncomingValues(ch chan utils.Message, config utils.Config) {

}
