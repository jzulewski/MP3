package nodes

import (
	"MP3/utils"
	"net"
	"sync"
)

func InitializeNodeServers(config *utils.Config) {
	for i, node := range config.Nodes {
		ln, err := net.Listen("tcp", ":"+node.Port)
		utils.CheckError(err)
		config.Nodes[i].Server = ln
	}
}

func InitializeNodeConnections(config *utils.Config) {
	nodes := config.Nodes
	ip := config.MServer.Ip
	port := config.MServer.Port
	CONNECT := ip + ":" + port

	// First connect each node to the master server
	// It will be the first connection in node.Conns
	for i := range nodes {

		//Connect to master server
		conn, err := net.Dial("tcp", CONNECT)
		utils.CheckError(err)

		// Append to actual struct
		nodes[i].Conns = append(nodes[i].Conns, conn)
	}

	// This node connects to the other node's server
	for i := range nodes {

		// To all these Nodes (including itself)
		for _, serverNode := range nodes {

			ip := serverNode.Ip
			port := serverNode.Port
			CONNECT := ip + ":" + port

			//Connect to server node
			conn, err := net.Dial("tcp", CONNECT)
			utils.CheckError(err)

			// Append to actual struct
			nodes[i].Conns = append(nodes[i].Conns, conn)
		}
	}
}

func StartSimulation(wg *sync.WaitGroup, config utils.Config) {
	for _, node := range config.Nodes {
		wg.Add(1)
		go handleNode(wg, node, config)
	}
}
