package server

import (
	"MP3/utils"
	"encoding/gob"
	"math"
	"net"
)

// High level function for the master server
func handleServer(config *utils.Config, ln net.Listener) {
	// Use channel to relay messages from receivers to processor
	ch := make(chan utils.Message)
	go handleConnections(ch, ln)
	processIncomingValues(ch, config)
}

// Accept all incoming connections
func handleConnections(ch chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go unicast_receive(ch, conn)
	}
}

// For each connection receive values and push them to the channel
func unicast_receive(ch chan utils.Message, conn net.Conn) {
	for {
		decoder := gob.NewDecoder(conn)
		var message utils.Message
		err := decoder.Decode(&message)
		utils.CheckError(err)
		ch <- message
	}
}

// Process incoming values
func processIncomingValues(ch chan utils.Message, config *utils.Config) {

	// Initialize variables
	n := len(config.Nodes)
	cf := 0
	states := make(map[string]float64)

	for {
		// For each message in the channel
		message := <-ch

		// If the node has failed, increment cf and delete node from states
		if message.Fail {
			cf++
			delete(states, message.From)
			continue
		} else {
			// Update node's state value and check if nodes should output
			states[message.From] = message.Value
			if len(states) == n-cf && checkStates(states) {
				// Multicast with output flag and break
				message := utils.Message{Output: true}
				multicast(message, config.MServer.Conns)
				break
			}
		}
	}
}

// Check if all values in states are within .001 of each other.
// Returns false if they are not within .001 of each other, true if they are.
func checkStates(states map[string]float64) bool {
	for _, value1 := range states {
		for _, value2 := range states {
			if math.Abs(value1-value2) > .001 {
				return false
			}
		}
	}
	return true
}

// Send message to each connection in conns without delay
func multicast(message utils.Message, conns []net.Conn) {
	for _, conn := range conns {
		encoder := gob.NewEncoder(conn)
		err := encoder.Encode(message)
		utils.CheckError(err)
	}
}
