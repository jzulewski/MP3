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
	for {
		var message utils.Message
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&message)
		utils.CheckError(err)
		valChan <- message.Value
	}

}

func unicast_send(nodeSend utils.Node, nodeReceive utils.Node) {

	ip := nodeReceive.Ip
	port := nodeReceive.Port
	CONNECT := ip + ":" + port

	//Connect to NodeReceive
	conn, err := net.Dial("tcp", CONNECT)
	utils.CheckError(err)

	nodeReceive.Conns = append(nodeReceive.Conns, conn)

	print(" Input: ")
	print(nodeSend.Input)
	print(" Id: ")
	print(nodeReceive.Id)

}

func processIncomingValues(valChan chan float64) {

}
