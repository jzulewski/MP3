# Distributed Systems MP3

# Overall System Design

The purpose of this exercise is to simulate an approximate consensus algorithm that relaxes agreement by requiring that nodes must output values within .001 of each other. Our design allows for N nodes, each with their own server. We also implemented a master server for keeping track of when the nodes should output, as well as other statistics useful in measuring the performance of the algorithm. In initialization, the master server and nodes start TCP listen servers. Then connections are established from every node to every other node, as well as from nodes to the master server and vice versa. These connections are saved in structs which can be accessed later. During the simulation, the node servers accept all incoming connections and push messages they receive to a function which processes those messages and determines what to do. A wait in main waits for all goroutines to finish before exiting the program.

# The Structs

* The Config struct stores data about the system. One instance of this struct is inialized and passed to functions which use what they need from it. It contains information from config.txt, including the delay parameters, the amount of nodes that can fail, and information about nodes and the master server.

* The Node struct contains information about each node. It holds the node's ID, input, IP, and Port specified in config.txt. The Conns field stores a list of connections that have been established to all other servers, with the connection to the master server being the first in the list. The Server field holds the listener of the server that was established for the node.

* The Server struct has information related to the master server. In addition to an IP and Port, it also has a list of connections that it has established with the nodes.

* The Message struct is how the nodes communicate. The Value and Round fields indicate the value and round being sent and the From field indicates who the message is from. The Output flag tells the nodes to output their current value. The Fail field means that the node that sent the message has failed.

# Packages

There are four packages: main, nodes, server, and utils. The main package contains only the main function which is run to start the program. The utils package contains utility functions like error checking and config struct initialization. The nodes and server packages are similar in that they both have init files. These files contain functions called by main to initialize servers, connections, and the simulation. In the nodes package, the node.go file is responsible for dealing with functions related to individual node processes like sending, receiving, and processing values. In the server package, the server.go file is responsible for dealing with functions related to master server processes like receiving messages from nodes and deciding when to output.

While this package 

# Sending Messages

Sometimes a node doesn't receive a message. In the `unicast_receive` function notice that a new gob decoder is created every time a message is decoded. This is intentional, as creating the decoder outside the for loop and reusing it results in an "extra data in buffer", as explained in [this github issue post](https://github.com/golang/go/issues/29766). However, this might create another issue, as sometimes two messages might be encoded at the same time, causing the decoder to decode only one of the messages and the other one is lost when a new decoder is created.

Another potential cause of message loss is on the encoding side. Every time the program sends a message a new encoder is created. This may cause issues as the encoder is designed to be able to be reused for the same stream. That being said, I haven't run into this problem in any previous MPs in which I created a new encoder each time a message was sent.

# Node Failure Implementation

Each node has a 3% chance of crashing each time it sends a message to another node. When a node fails it sends a message to every other node with the Fail flag set to true. Each node keeps track of how many nodes have failed and the maximum number that are allowed to fail. When a node receives a message with the Fail flag, it increments the current failed nodes and checks if it's greater than or equal to f. If it is, then the node knows it cannot crash and sets the variable canFail (initialized as true) to false. This bypasses the 3% chance to fail the next time multicast is called.

If nodes crash at the same time it's possible that more than f nodes will crash. This is a limitation of the speed of sending messages to concurrently running processes over TCP, which can't be instant. Additionally, because rand is seeded with nano unix time, sometimes goroutines reach the seed line at the same time and end up with the same seed. This results in multiple nodes crashing at the same time. Ideally a better source of randomness would be used.

# Storage of Future Values

When a node receives a message with a value of a round greater than the round it's currently on it needs to store that value somehow. Currently the message is put back into the channel. This could result in an infinite loop until a new message with a round equal to the node's round is pushed to the channel. A better implemtation might involved storing the message in a list, but that would take more memory.

repeated functions handleConnections and unicast_receive in server and node packages.
