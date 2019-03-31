### ElevatorSync

This module is responsible for the communication between orderManager and the network module. The module exchanges information about the local and global elevators. 

## Awknowledgement

As the elevators use a peer-to-peer architecture, it is crucial that all elevators have the same information on eachother. This will make sure no orders are lost in case of motorfailure or power-loss on one of the computers, as well as ensuring the hallAssigner will assign an order to the correct elevator. This task is done by using an awknowledgement matrix, which looks like the elevator matrix from orderManager, but that has every index replaced with an AckStruct. The AckStruct consists of the data sent, and an array of bool, verifying if an elevator is awaiting an awknowledgement. By periodically check if the matrix has any awaitingAck set to true, this module make sure everyone has the same information, which also acts as an error handler in case of packet-loss.

