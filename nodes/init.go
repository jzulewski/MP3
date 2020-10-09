package nodes

import (
	"MP3/utils"
	"net"
)

func InitializeNodes(nodes []utils.Node) {
	for _, node := range nodes {
		ln := startServer(node)
		go handleNode(node, ln)
	}
}

func startServer(node utils.Node) net.Listener {
	ln, err := net.Listen("tcp", ":"+node.Port)
	utils.CheckError(err)
	return ln
}

func InitializeConnections() {

}
