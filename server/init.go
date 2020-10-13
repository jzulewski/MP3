package server

import (
	"MP3/utils"
	"net"
)

func InitializeMasterServer(config *utils.Config) {
	ln, err := net.Listen("tcp", ":"+config.MServer.Port)
	utils.CheckError(err)
	go handleServer(config, ln)
}

func InitializeServerConnections(config *utils.Config) {
	nodes := config.Nodes
	for _, node := range nodes {

		ip := node.Ip
		port := node.Port
		CONNECT := ip + ":" + port

		//Connect to node's server
		conn, err := net.Dial("tcp", CONNECT)
		utils.CheckError(err)

		// Append to server conns list
		config.MServer.Conns = append(config.MServer.Conns, conn)
	}
}
