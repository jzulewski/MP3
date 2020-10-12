package server

import (
	"MP3/utils"
	"net"
)

func InitializeMasterServer(config utils.Config) {
	ln, err := net.Listen("tcp", ":"+config.MServer.Port)
	utils.CheckError(err)
	handleServer(config, ln)
}

func InitializeServerConnections(server utils.Server, nodes []utils.Node) {
	for _, node := range nodes {

		ip := node.Ip
		port := node.Port
		CONNECT := ip + ":" + port

		//Connect to node's server
		conn, err := net.Dial("tcp", CONNECT)
		utils.CheckError(err)

		// Append to server conns list
		server.Conns = append(server.Conns, conn)
	}
}
