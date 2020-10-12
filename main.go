package main

import (
	"MP3/nodes"
	"MP3/utils"
)

func main() {
	// Initialize config struct
	config := utils.CreateConfigStruct()

	// Initialize node processes
	nodes.InitializeNodeServers(config)

	// Initialize connections between nodes
	nodes.InitializeConnections(config.Nodes)

	// Begin Simulation
	nodes.StartSimulation(config)

	ch := make(chan int)
	temp := <-ch
	print(temp)
}
