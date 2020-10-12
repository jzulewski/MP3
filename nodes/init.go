package nodes

import (
	"MP3/utils"
	"net"
)

func InitializeNodeServers(config utils.Config) {
	for i, node := range config.Nodes {
		ln, err := net.Listen("tcp", ":"+node.Port)
		utils.CheckError(err)
		config.Nodes[i].Server = ln
	}
}

func InitializeConnections(nodes []utils.Node) {

	//this Node connects to the server
	for i := range nodes {

		//to all these Nodes (including itself)
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

func StartSimulation(config utils.Config) {
	for _, node := range config.Nodes {
		go handleNode(node, config)
	}
}
