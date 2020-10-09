package nodes

import (
	"MP3/utils"
	"encoding/gob"
	"net"
)

func handleNode(node utils.Node, ln net.Listener) {
	valChan := make(chan float64)
	go handleConnections(valChan, node, ln)
	processIncomingValues(valChan)
}

func handleConnections(valChan chan float64, node utils.Node, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go unicast_receive(valChan, conn)
	}
}

func unicast_receive(valChan chan float64, conn net.Conn) {
	var message utils.Message
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&message)
	utils.CheckError(err)
	valChan <- message.Value
}

func processIncomingValues(valChan chan float64) {

}
