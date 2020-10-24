package nodes

import (
	"MP3/utils"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// Upper level function of each node. Links message reception with processing
func handleNode(wg *sync.WaitGroup, node utils.Node, config utils.Config) {
	defer wg.Done()
	// Use channel to communicate messages from receivers to processors
	valChan := make(chan utils.Message)
	go handleConnections(valChan, node.Server)
	processIncomingValues(valChan, node, config)
}

// Accept each incoming connection on its own goroutine
func handleConnections(valChan chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go unicast_receive(valChan, conn)
	}
}

// Receive messages on the connection and push them to the channel
func unicast_receive(valChan chan utils.Message, conn net.Conn) {
	for {
		decoder := gob.NewDecoder(conn)
		var message utils.Message
		err := decoder.Decode(&message)
		utils.CheckError(err)
		valChan <- message
	}
}

// Process messages from the valChan channel
func processIncomingValues(valChan chan utils.Message, node utils.Node, config utils.Config) {

	// Variable initialization
	from := node.Id
	n := len(config.Nodes)
	f := config.F
	cf := 0
	round := 1
	var value float64
	recievedMessages := 0
	sum := 0.
	canFail := true

	// Initial multicast of the node's input value
	err := multicast(node.Conns,
		utils.Message{From: from, Round: round, Value: node.Input},
		true,
		config.MinDelay,
		config.MaxDelay)
	// If the node crashed during multicast, handle it
	if err != nil {
		fmt.Println("Node "+node.Id+":", err.Error())
		// Multicast to other nodes that this node has crashed
		err = multicast(node.Conns,
			utils.Message{From: from, Fail: true},
			false,
			0, 1)
		utils.CheckError(err)
		return
	}

	// Infinite for loop
	for {
		// For each message in the channel
		message := <-valChan

		// Check flags and round number
		if message.Output {
			println(value)
			break
		} else if message.Fail {
			cf++
			if cf >= f {
				canFail = false
			}
		} else if message.Round < round {
			continue
		} else if message.Round > round {
			valChan <- message
		} else {
			// Message is for current round
			recievedMessages++
			sum += message.Value

			// If node has received n-f messages this round, update value and multicast
			if recievedMessages >= n-f {
				value = sum / float64(n-f)
				fmt.Printf("Node %s finished round %d. Value %f\n", from, round, value)
				err = multicast(node.Conns,
					utils.Message{From: from, Round: round + 1, Value: value},
					canFail,
					config.MinDelay,
					config.MaxDelay)
				// If node crashed during multicast, send message to all other nodes and master server
				if err != nil {
					fmt.Println("Node "+node.Id+":", err.Error())
					err = multicast(node.Conns,
						utils.Message{From: from, Fail: true},
						false,
						0, 1)
					utils.CheckError(err)
					return
				}
				round++
				recievedMessages = 0
				sum = 0.
			}
		}
	}
}

// Send message to every connection in conns, with min and max delay in ms
func multicast(conns []net.Conn, message utils.Message, canFail bool, min, max int) error {

	// Always send to master server
	encoder := gob.NewEncoder(conns[0])
	err := encoder.Encode(message)
	utils.CheckError(err)

	// Send to all other conns with a chance to crash
	rand.Seed(time.Now().UnixNano())
	for _, conn := range conns[1:] {
		if canFail {
			number := rand.Intn(100)
			if number < 3 {
				return errors.New("node has crashed")
			}
		}
		// Send as goroutine to avoid bottleneck
		go unicast_send(conn, message, min, max)
	}

	return nil
}

// Send message to designated connection with delay
func unicast_send(conn net.Conn, message utils.Message, min, max int) {
	delay := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(message)
	utils.CheckError(err)
}
