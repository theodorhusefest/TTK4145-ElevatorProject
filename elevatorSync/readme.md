### ElevatorSync

This module is responsible for the communication between orderManager and the network module. The module exchanges information about the local and global elevators.

## Awknowledgement

This module is also responsible for error handling is case of packet loss. This is done by using an awknowledgement matrix, which looks like the elevator matrix from orderManager, but that has every index replaced with an AckStruct. The AckStruct consists of the data sent, and an array of bool, verifying if an elevator is awaiting an awknowledgement.
