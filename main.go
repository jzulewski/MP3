package main

import (
	"MP3/nodes"
	"MP3/utils"
)

func main() {
	// Initialize config struct
	config := utils.CreateConfigStruct()

	// Initialize node processes
	nodes.InitializeNodes(config.Nodes)

	// Initialize connections between nodes
	nodes.InitializeConnections()
}
