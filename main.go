package main

import (
	"MP3/nodes"
	"MP3/server"
	"MP3/utils"
	"sync"
)

func main() {
	// Initialize config struct
	config := utils.CreateConfigStruct()

	// Initialize master server
	server.InitializeMasterServer(&config)

	// Initialize node processes
	nodes.InitializeNodeServers(&config)

	// Initialize connections between master server and node servers
	server.InitializeServerConnections(&config)

	// Initialize connections between nodes
	nodes.InitializeNodeConnections(&config)

	// Begin Simulation
	var wg sync.WaitGroup
	nodes.StartSimulation(&wg, config)

	// Wait for goroutines to finish
	wg.Wait()
	println("Simulation finished")
}
