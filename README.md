# Distributed Systems MP3

# Overall System Design

The purpose of this exercise is to simulate an approximate consensus algorithm that relaxes agreement by requiring that nodes must output values within .001 of each other. Our design allows for N nodes, each with their own server. We also implemented a master server for keeping track of when the nodes should output, as well as other statistics useful in measuring the performance of the algorithm. In initialization, the master server and nodes start TCP listen servers. Then connections are established from every node to every other node, as well as from nodes to the master server and vice versa. These connections are saved in structs which can be accessed later. During the simulation, the node servers accept all incoming connections and push messages they receive to a function which processes those messages and determines what to do. A wait in main waits for all goroutines to finish before exiting the program.

# The Structs

* The Config struct stores data about the system. One instance of this struct is inialized and passed to functions which use what they need from it. It contains information from config.txt, including the delay parameters, the amount of nodes that can fail, and information about nodes and the master server.

* The Node struct contains information about each node. It holds the node's ID, input, IP, and Port specified in config.txt. The Conns field stores a list of connections that have been established to all other servers, with the connection to the master server being the first in the list. The Server field holds the listener of the server that was established for the node.

* The Server struct has information related to the master server. In addition to an IP and Port, it also has a list of connections that it has established with the nodes.

* The Message struct is how the nodes communicate. The Value and Round fields indicate the value and round being sent and the From field indicates who the message is from. The Output flag tells the nodes to output their current value. The Fail field means that the node that sent the message has failed.

use of future messages list does not always keep order of incoming messages which might impact something.

repeated functions handleConnections and unicast_receive in server and node packages.
