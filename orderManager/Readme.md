### Order Manager

This is the module responsible to have control over the elevator matrix, forwarding information to the sync-module from FSM and triggering FSM when having orders to be done.
The elevator matrix for a single elevator is structured like this:

|	  			| 			| 			| 
| --   			| --			| --			|
| ID   			|   		|   		| 
| State			|  			| 			| 
| Floor			|  			| 			| 
| Direction		|  			| 			| 
| ----------- 	| Down 4 	| Cab 4		|
| Up 3 			| Down 3 	| Cab 3 	|
| Up 2 			| Down 2 	| Cab 2		|
| Up 1 			| --------- | Cab 1		|

Where every value is stored as an integer. ID is stored as values 0, 1, 2, ..., n-1 elevators. Floor is stored as 1, 2, ..., n-1 floors. Direction is stored as -1, 0, or 1. The rest of the matrix values are stored as simple ones or zeros, as these represent orders pushed into the button board.

## Connection loss
In order to not loose any information due to a connection loss, each elevator stores the matrix of both its own and the other elevators' matrices locally. In this way, should connection to an elevator be lost, the information can simply be transmitted to from one of the other elevators once the lost elevator comes back online. The matrix then looks something like this:


| elevator matrix 1 | elevator matrix 2 | elevator matrix 3 | ... | elevator matrix n
| --  | -- | -- | -- | -- |

