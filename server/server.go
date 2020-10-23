package server

import (
	"MP3/utils"
	"encoding/gob"
	"math"
	"net"
)

func handleServer(config *utils.Config, ln net.Listener) {
	ch := make(chan utils.Message)
	go handleConnections(ch, ln)
	processIncomingValues(ch, config)
}

func handleConnections(ch chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go unicast_receive(ch, conn)
	}
}

func unicast_receive(ch chan utils.Message, conn net.Conn) {
	for {
		decoder := gob.NewDecoder(conn)
		var message utils.Message
		err := decoder.Decode(&message)
		utils.CheckError(err)
		ch <- message
	}
}

func processIncomingValues(ch chan utils.Message, config *utils.Config) {
	n := len(config.Nodes)
	cf := 0
	states := make(map[string]float64)
	for {
		message := <-ch
		if message.Fail {
			cf++
			delete(states, message.From)
			continue
		} else {
			states[message.From] = message.Value
			if len(states) == n-cf && checkStates(states) {
				message := utils.Message{Output: true}
				multicast(message, config.MServer.Conns)
			}
		}
	}
}

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

func multicast(message utils.Message, conns []net.Conn) {
	for _, conn := range conns {
		encoder := gob.NewEncoder(conn)
		err := encoder.Encode(message)
		utils.CheckError(err)
	}
}
