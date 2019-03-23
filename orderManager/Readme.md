## Order Manager

This is the module responsible to have control over the elevator matrix, forwarding information to the sync-module from FSM and triggering FSM when having orders to be done.
The elevator matrix for a single elevator is structured like this: 
|	  			| 			| 			| 
|--   			|--			|--			|
|ID   			|   		|   		| 
|State			|  			| 			| 
|Floor			|  			| 			| 
|Direction		|  			| 			| 
|----------- 	| Down 4 	| Cab 4		|
|Up 3 			| Down 3 	| Cab 3 	|
|Up 2 			| Down 2 	| Cab 2		|
|Up 1 			| --------- | Cab 1		|

