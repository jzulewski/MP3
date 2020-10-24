# Distributed Systems MP3

# To Run

cd to the MP3 directory and type
`go build MP3`

# Overall System Design

The purpose of this exercise is to simulate an approximate consensus algorithm that relaxes agreement by requiring that nodes must output values within .001 of each other. Our design allows for N nodes, each with their own server. We also implemented a master server for keeping track of when the nodes should output, as well as other statistics useful in measuring the performance of the algorithm. In initialization, the master server and nodes start TCP listen servers. Then connections are established from every node to every other node, as well as from nodes to the master server and vice versa. These connections are saved in structs which can be accessed later. During the simulation, the node servers accept all incoming connections and push messages they receive to a function which processes those messages and determines what to do. A wait in main waits for all goroutines to finish before exiting the program.

# Config.txt

The first line of config.txt is the minimum delay (ms), the maximum delay (ms), and the maximum number of node failures separated by a space. The second line is the IP and Port of the master server separated by a space. The subsequent lines are nodes. The format for nodes is `ID input IP Port`. An example of a valid config.txt is below.
```
1000 5000 1
127.0.0.1 5000
1 0.3 127.0.0.1 5001
2 0.4 127.0.0.1 5002
3 0.5 127.0.0.1 5003
4 0.9 127.0.0.1 5004
```

# The Structs

* The Config struct stores data about the system. One instance of this struct is inialized and passed to functions which use what they need from it. It contains information from config.txt, including the delay parameters, the amount of nodes that can fail, and information about nodes and the master server.

* The Node struct contains information about each node. It holds the node's ID, input, IP, and Port specified in config.txt. The Conns field stores a list of connections that have been established to all other servers, with the connection to the master server being the first in the list. The Server field holds the listener of the server that was established for the node.

* The Server struct has information related to the master server. In addition to an IP and Port, it also has a list of connections that it has established with the nodes.

* The Message struct is how the nodes communicate. The Value and Round fields indicate the value and round being sent and the From field indicates who the message is from. The Output flag tells the nodes to output their current value. The Fail field means that the node that sent the message has failed.

# Packages

There are four packages: main, nodes, server, and utils. The main package contains only the main function which is run to start the program. The utils package contains utility functions like error checking and config struct initialization. The nodes and server packages are similar in that they both have init files. These files contain functions called by main to initialize servers, connections, and the simulation. In the nodes package, the node.go file is responsible for dealing with functions related to individual node processes like sending, receiving, and processing values. In the server package, the server.go file is responsible for dealing with functions related to master server processes like receiving messages from nodes and deciding when to output.

While this package layout makes sense, some code is repeated in the server and nodes package which isn't ideal. Specifically, both nodes and the master server need to accept connections to the server, decode messages on those connections, and send the message to be processed. This design is encapsulated in the `handleServer` and `handleNode` functions, which are nearly identical pieces of code, repeated in different packages. Ideally this code block would only be written once and would called by the packages when necessary.

# Sending Messages

Sometimes a node doesn't receive a message. In the `unicast_receive` function notice that a new gob decoder is created every time a message is decoded. This is intentional, as creating the decoder outside the for loop and reusing it results in an "extra data in buffer", as explained in [this github issue post](https://github.com/golang/go/issues/29766). However, this might create another issue, as sometimes two messages might be encoded at the same time, causing the decoder to decode only one of the messages and the other one is lost when a new decoder is created.

Another potential cause of message loss is on the encoding side. Every time the program sends a message a new encoder is created. This may cause issues as the encoder is designed to be able to be reused for the same stream. That being said, I haven't run into this problem in any previous MPs in which I created a new encoder each time a message was sent.

In truth, nothing tangible could be found online to fix this problem, and even the good people of /r/golang were not able to help. The result of this bug is that sometimes the program will "hang", because if a node does not receive n-f values in a round it will forever be waiting for more messages that will never come. This causes a chain reaction as other nodes, expecting the frozen node to send a message, will also freeze up. Additionally, sometimes a crashed node's message to the master server will fail, causing the simulation to run infinitely even when it should have stopped. Luckily, due to the randomness of the bug's appearance, simulations were still able to be run and data was collected from those successful simulation runs.

# Node Failure Implementation

Each node has a 3% chance of crashing each time it sends a message to another node. When a node fails it sends a message to every other node with the Fail flag set to true. Each node keeps track of how many nodes have failed and the maximum number that are allowed to fail. When a node receives a message with the Fail flag, it increments the current failed nodes and checks if it's greater than or equal to f. If it is, then the node knows it cannot crash and sets the variable canFail (initialized as true) to false. This bypasses the 3% chance to fail the next time multicast is called.

If nodes crash at the same time it's possible that more than f nodes will crash. This is a limitation of the speed of sending messages to concurrently running processes over TCP, which can't be instant. Additionally, because rand is seeded with nano unix time, sometimes goroutines reach the seed line at the same time and end up with the same seed. This results in multiple nodes crashing at the same time. Ideally a better source of randomness would be used.

# When to Output

To decide when to output, we used a master server that recieves and analyzes the values of the nodes. The master server stores these values in a map that maps the node ids to their float64 value. Each time the map is updated, the state of each node is checked against the state of every other node, and if the values are all within .001 of each other, the server sends a message to all nodes with the Output flag set to true. When nodes receive a message with the output flag they print their value and return from the function, effectively "closing" that node. When a node fails, the server makes sure to remove it from the map, because it's value no longer matters in determining the agreement property.

One benefit to the master server approach to determining agreement is that the server has full information. While nodes only receive n-f messages per round, the master server receives all messages from non-crashed nodes. This means that it is never wrong about when to ouput, whereas a node working with partial information may think it should output when in reality it should not. One downside to this approach is that it is not very representative of reality because a master server with full information that never crashes is not very realistic.

# Storage of Future Values

When a node receives a message with a value of a round greater than the round it's currently on it needs to store that value somehow. Currently the message is put back into the channel. This could result in an infinite loop until a new message with a round equal to the node's round is pushed to the channel. A better implemtation might involved storing the message in a list, but that would take more memory.
