package nodes

import (
	"MP3/utils"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

func handleNode(wg *sync.WaitGroup, node utils.Node, config utils.Config) {
	defer wg.Done()
	valChan := make(chan utils.Message)
	go handleConnections(valChan, node.Server)
	processIncomingValues(valChan, node, config)
}

func handleConnections(valChan chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go unicast_receive(valChan, conn)
	}
}

func unicast_receive(valChan chan utils.Message, conn net.Conn) {
	for {
		var message utils.Message
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&message)
		utils.CheckError(err)
		valChan <- message
	}
}

func processIncomingValues(valChan chan utils.Message, node utils.Node, config utils.Config) {

	from := node.Id
	n := len(config.Nodes)
	f := config.F

	round := 1
	var value float64
	recievedMessages := 0
	sum := 0.
	var futureMessages []utils.Message

	go multicast(node.Conns,
		utils.Message{From: from, Round: round, Value: node.Input},
		config.MinDelay,
		config.MaxDelay)

	for {
		message := <-valChan

		if message.Output {
			println(value)
			break
		} else if message.Round < round {
			continue
		} else if message.Round == round {
			recievedMessages++
			sum += message.Value
		} else if message.Round > round {
			futureMessages = append(futureMessages, message)
		}

		if recievedMessages >= n-f {
			value = sum / float64(n-f)
			fmt.Printf("Node %s, Value %f, Round %d\n", from, value, round)
			go multicast(node.Conns,
				utils.Message{From: from, Round: round + 1, Value: value},
				config.MinDelay,
				config.MaxDelay)
			for _, message := range futureMessages {
				valChan <- message
			}
			round++
			recievedMessages = 0
			sum = 0.
		}
	}
}

func multicast(conns []net.Conn, message utils.Message, min, max int) {

	// Always send to master server
	encoder := gob.NewEncoder(conns[0])
	err := encoder.Encode(message)
	utils.CheckError(err)

	// Send to all other conns
	//rand.Seed(time.Now().UnixNano())
	for _, conn := range conns[1:] {
		//number := rand.Intn(100)
		//if number < 95 {
		//	unicast_send(conn, message, min, max)
		//}
		go unicast_send(conn, message, min, max)
	}
}

func unicast_send(conn net.Conn, message utils.Message, min, max int) {
	delay := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(message)
	utils.CheckError(err)
}
