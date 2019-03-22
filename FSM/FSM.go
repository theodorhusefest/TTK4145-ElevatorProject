package FSM

import (
	. "../Config"
	"../IO"
	//"../Utilities"
	"../orderManager"
	"fmt"
	"time"
)

type FSMchannels struct {
	NewLocalOrderChan  chan int
	ArrivedAtFloorChan chan int
	DoorTimeoutChan    chan bool
}

func StateMachine(FSMchans FSMchannels, LocalOrderFinishedChan chan int, UpdateElevStatusch chan Message, elevatorMatrix [][]int, elevator Elevator) {

	doorOpenTimeOut := time.NewTimer(3 * time.Second)
	motorFailureTimeOut := time.NewTimer(5 * time.Second)
	motorFailureTimeOut.Stop()
	doorOpenTimeOut.Stop()

	for {
		select {
		case newLocalOrder := <-FSMchans.NewLocalOrderChan:

			switch elevator.State {
			case IDLE:

				elevator.Dir = chooseDirection(elevatorMatrix, elevator)
				io.SetMotorDirection(elevator.Dir)
				motorFailureTimeOut.Reset(5 * time.Second)

				orderManager.InsertDirection(elevator.ID, elevator.Dir, elevatorMatrix)

				if elevator.Dir == DIR_Stop {
					io.SetDoorOpenLamp(true)
					elevator.State = DOOROPEN
					orderManager.InsertState(elevator.ID, int(DOOROPEN), elevatorMatrix)
					doorOpenTimeOut.Reset(3 * time.Second)
					LocalOrderFinishedChan <- elevator.Floor
				}
				elevator.State = MOVING
				orderManager.InsertState(elevator.ID, int(MOVING), elevatorMatrix)

			case MOVING:

			case DOOROPEN:
				elevator.Dir = chooseDirection(elevatorMatrix, elevator)
				if elevator.Floor == newLocalOrder {
					doorOpenTimeOut.Reset(3 * time.Second)
					fmt.Println("Resetting Time")
					LocalOrderFinishedChan <- elevator.Floor
				}
			}

			updatedStates := Message{Select: UpdateStates, ID: elevator.ID, State: int(elevator.State), Floor: elevator.Floor, Dir: elevator.Dir}
			UpdateElevStatusch <- updatedStates

		case currentFloor := <-FSMchans.ArrivedAtFloorChan:
			orderManager.InsertFloor(elevator.ID, currentFloor, elevatorMatrix)
			elevator.Floor = currentFloor
			io.SetFloorIndicator(currentFloor)


			if shouldStop(elevator.ID, elevator, elevatorMatrix) {
				elevator.State = DOOROPEN
				io.SetDoorOpenLamp(true)
				doorOpenTimeOut.Reset(3 * time.Second)
				io.SetMotorDirection(DIR_Stop)

				orderManager.InsertState(elevator.ID, int(DOOROPEN), elevatorMatrix)
				orderManager.InsertDirection(elevator.ID, elevator.Dir, elevatorMatrix)

				LocalOrderFinishedChan <- elevator.Floor
			}
			updatedStates := Message{Select: UpdateStates, ID: elevator.ID, State: int(elevator.State), Floor: elevator.Floor, Dir: elevator.Dir}
			UpdateElevStatusch <- updatedStates

		case <-doorOpenTimeOut.C:
			fmt.Println("DOOROPEN, Current Direction is ", elevator.Dir)
			io.SetDoorOpenLamp(false)
			elevator.Dir = chooseDirection(elevatorMatrix, elevator)
			orderManager.InsertDirection(elevator.ID, elevator.Dir, elevatorMatrix)
			io.SetMotorDirection(elevator.Dir)
			LocalOrderFinishedChan <- elevator.Floor
			if elevator.Dir == DIR_Stop {
				elevator.State = IDLE
				orderManager.InsertState(elevator.ID, int(IDLE), elevatorMatrix)
			} else {
				//io.SetMotorDirection(elevator.Dir)
				elevator.State = MOVING
				orderManager.InsertState(elevator.ID, int(MOVING), elevatorMatrix)

				motorFailureTimeOut.Reset(5 * time.Second)
			}
			updatedStates := Message{Select: UpdateStates, ID: elevator.ID, State: int(elevator.State), Floor: elevator.Floor, Dir: elevator.Dir}
			UpdateElevStatusch <- updatedStates
		}
	}
}

func chooseDirection(elevatorMatrix [][]int, elevator Elevator) MotorDirection {

	switch elevator.Dir {
	case DIR_Up:
		if isOrderAbove(elevator.ID, elevator.Floor, elevatorMatrix) {
			return DIR_Up
		}
		if isOrderBelow(elevator.ID, elevator.Floor, elevatorMatrix) {
			return DIR_Down
		}
		return DIR_Stop
	
	case DIR_Down:
		if isOrderBelow(elevator.ID, elevator.Floor, elevatorMatrix) {
			return DIR_Down
		}
		if isOrderAbove(elevator.ID, elevator.Floor, elevatorMatrix) {
			return DIR_Up
		}
		return DIR_Stop

	default: 
		if isOrderAbove(elevator.ID, elevator.Floor, elevatorMatrix) {
			return DIR_Up
		}
		if isOrderBelow(elevator.ID, elevator.Floor, elevatorMatrix) {
			return DIR_Down
		}
		return DIR_Stop


	}


}

func isOrderAbove(id int, currentFloor int, elevatorMatrix [][]int) bool {
	if currentFloor == 3 {
		return false
	}

	for floor := (len(elevatorMatrix) - currentFloor - 2); floor > 3; floor-- {
		for buttons := (id * 3); buttons < (id*3 + 3); buttons++ {
			if elevatorMatrix[floor][buttons] == 1 {
				return true
			}
		}
	}
	return false
}

func isOrderBelow(id int, currentFloor int, elevatorMatrix [][]int) bool {
	if currentFloor == 0 {
		return false
	}

	for floor := (len(elevatorMatrix) - currentFloor); floor < (len(elevatorMatrix)); floor++ {
		for buttons := (id * 3); buttons < (id*3 + 3); buttons++ {
			if elevatorMatrix[floor][buttons] == 1 {
				return true
			}
		}

	}
	return false
}

func shouldStop(id int, elevator Elevator, elevatorMatrix [][]int) bool {
	//Cab call is pressed, stop
	if elevatorMatrix[len(elevatorMatrix)-elevator.Floor-1][id*NumElevators+2] == 1 {
		return true
	}
	switch elevator.Dir {
	case DIR_Up:
		if elevatorMatrix[len(elevatorMatrix)-elevator.Floor-1][id*NumElevators] == 1 {
			return true
		} else if elevatorMatrix[len(elevatorMatrix)-elevator.Floor-1][id*NumElevators+1] == 1 && !isOrderAbove(id, elevator.Floor, elevatorMatrix) {
			return true
		}
	case DIR_Down:
		if elevatorMatrix[len(elevatorMatrix)-elevator.Floor-1][id*NumElevators+1] == 1 {
			return true
		} else if elevatorMatrix[len(elevatorMatrix)-elevator.Floor-1][id*NumElevators] == 1 && !isOrderBelow(id, elevator.Floor, elevatorMatrix) {
			return true
		}
	}
	return false
}
