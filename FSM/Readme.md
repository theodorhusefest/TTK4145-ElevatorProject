## FSM

This is the module which controls the elevator, and does so by switching between the four states:
 - IDLE
 - MOVING 
 - DOOROPEN
 - UNDEFINED

This module will trigger a change in the elevator-state when recieving an event from orderManager, and is closely connected to the IO-module. The FSM is responsible to execute all orders in the elevator matrix, and detects if there is motor failure. 


