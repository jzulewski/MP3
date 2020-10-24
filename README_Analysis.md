**Analysis of MP3**

Overall, we noticed two notable trends. First, that nodes having some initial consensus improved performance (which is measured by time and number of rounds)
and secondly, there is not any clear correlation, positive or negative, between the number of nodes and performance.

**Initial Consensus**
"Initial Consensus" is what we call having some nodes with the same initial inputs. In this case, we had 8 nodes with initial consensus (Two nodes each with .3, .4, .5, .9 values) and the nodes without initial consensus had random values for their inputs. To limit confounding variables, the two sets of nodes had the same number of failures, delays, and range of their inputs. Based on our observation of the data, it is clear that initial consensus improves performance. We believe this is because the n-f messages that each node receives and calculates into their output are more likely to be the same as another node since there are repeating values among the nodes, and therefore consensus will be reacher quicker.
![Screenshot](Screen%20Shot%202020-10-23%20at%209.57.20%20PM.png)

**Number of Nodes**
![Screenshot](Screen%20Shot%202020-10-23%20at%209.47.57%20PM.png)
In this experiment, we had one group of 4 nodes and another of 8 nodes. Other factors such as input range, number of failures, and delays were kept the same between the groups. Based on our observations, we noticed no correlation between the number of nodes and performance. As seen in the middle of the graph, sometimes the 4 nodes group was faster and other times it was the 8 nodes groups. Additionally, outliers occured such as when the first node crashed and the number of failures was limited to 1, which lead to all the rest of the nodes receiving the same 7 inputs in the first round and therefore reaching consensus usually quickly. Another outlier appears to be when the 4 nodes took 4 rounds and potentially could have had longer delays. Overall, we can conclude that an increase or decrease in the number of nodes does not have a direct relationship to performance.

**Other Trends**
In general, input ranges did not have a direct relationship to performance
Increase Delay Times (both the minimum or maximum) worsened performance
Our code had some bugs for crashing nodes, but throughout our testing, we did not notice the number of failures having a significicant impact on performance


