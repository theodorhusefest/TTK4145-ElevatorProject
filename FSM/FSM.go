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

func StateMachine(	elevatorMatrix [][]int, localElev Elevator, 
					FSMchans FSMchannels, LocalOrderFinishedChan chan int, 
					UpdateElevStatusch chan Message) {

	doorOpenTimeOut := time.NewTimer(3 * time.Second)
	motorFailureTimeOut := time.NewTimer(5 * time.Second)
	motorFailureTimeOut.Stop()
	doorOpenTimeOut.Stop()

	for {
		select {
		case newLocalOrder := <-FSMchans.NewLocalOrderChan:



			switch localElev.State {
			case IDLE:

				localElev.Dir = chooseDirection(elevatorMatrix, localElev)
				io.SetMotorDirection(localElev.Dir)

				orderManager.InsertDirection(localElev.ID, localElev.Dir, elevatorMatrix)

				if localElev.Dir == DIR_Stop {
					io.SetDoorOpenLamp(true)
					localElev.State = DOOROPEN
					orderManager.InsertState(localElev.ID, int(DOOROPEN), elevatorMatrix)
					doorOpenTimeOut.Reset(3 * time.Second)
					LocalOrderFinishedChan <- localElev.Floor
				} else {
					localElev.State = MOVING
					orderManager.InsertState(localElev.ID, int(MOVING), elevatorMatrix)
					motorFailureTimeOut.Reset(5 * time.Second)
				}


			case MOVING:

			case DOOROPEN:
				localElev.Dir = chooseDirection(elevatorMatrix, localElev)
				if localElev.Floor == newLocalOrder {
					doorOpenTimeOut.Reset(3 * time.Second)
					LocalOrderFinishedChan <- localElev.Floor
				}

			case UNDEFINED:
				fmt.Println("Motor has failed")

			}
			updatedStates := Message{Select: UpdateStates, ID: localElev.ID, State: int(localElev.State), Floor: localElev.Floor, Dir: localElev.Dir}
			UpdateElevStatusch <- updatedStates

		case currentFloor := <-FSMchans.ArrivedAtFloorChan:

			orderManager.InsertFloor(localElev.ID, currentFloor, elevatorMatrix)
			localElev.Floor = currentFloor
			io.SetFloorIndicator(currentFloor)

			if shouldStop(localElev.ID, localElev, elevatorMatrix) {
				localElev.State = DOOROPEN
				io.SetDoorOpenLamp(true)
				io.SetMotorDirection(DIR_Stop)
				doorOpenTimeOut.Reset(3 * time.Second)
				motorFailureTimeOut.Stop()

				orderManager.InsertState(localElev.ID, int(DOOROPEN), elevatorMatrix)
				orderManager.InsertDirection(localElev.ID, localElev.Dir, elevatorMatrix)

				LocalOrderFinishedChan <- localElev.Floor
			} else if localElev.State != IDLE {
				motorFailureTimeOut.Reset(5 * time.Second)
			}
			updatedStates := Message{Select: UpdateStates, ID: localElev.ID, State: int(localElev.State), Floor: localElev.Floor, Dir: localElev.Dir}
			UpdateElevStatusch <- updatedStates

		case <-doorOpenTimeOut.C:
			io.SetDoorOpenLamp(false)
			localElev.Dir = chooseDirection(elevatorMatrix, localElev)
			orderManager.InsertDirection(localElev.ID, localElev.Dir, elevatorMatrix)
			io.SetMotorDirection(localElev.Dir)
			LocalOrderFinishedChan <- localElev.Floor
			if localElev.Dir == DIR_Stop {
				localElev.State = IDLE
				orderManager.InsertState(localElev.ID, int(IDLE), elevatorMatrix)
				motorFailureTimeOut.Stop()

			} else {
				localElev.State = MOVING
				orderManager.InsertState(localElev.ID, int(MOVING), elevatorMatrix)
				motorFailureTimeOut.Reset(5 * time.Second)
			}

			updatedStates := Message{Select: UpdateStates, ID: localElev.ID, State: int(localElev.State), Floor: localElev.Floor, Dir: localElev.Dir}
			UpdateElevStatusch <- updatedStates

		case <- motorFailureTimeOut.C:
			fmt.Println("Motor has failed")
			localElev.State = UNDEFINED
			orderManager.InsertState(localElev.ID, int(UNDEFINED), elevatorMatrix)

			updatedStates := Message{Select: UpdateStates, ID: localElev.ID, State: int(localElev.State), Floor: localElev.Floor, Dir: localElev.Dir}
			UpdateElevStatusch <- updatedStates
		}
	}
}

func chooseDirection(elevatorMatrix [][]int, localElev Elevator) MotorDirection {

	switch localElev.Dir {
	case DIR_Up:
		if isOrderAbove(localElev.ID, localElev.Floor, elevatorMatrix) {
			return DIR_Up
		}
		if isOrderBelow(localElev.ID, localElev.Floor, elevatorMatrix) {
			return DIR_Down
		}
		return DIR_Stop

	case DIR_Down:
		if isOrderBelow(localElev.ID, localElev.Floor, elevatorMatrix) {
			return DIR_Down
		}
		if isOrderAbove(localElev.ID, localElev.Floor, elevatorMatrix) {
			return DIR_Up
		}
		return DIR_Stop

	default:
		if isOrderAbove(localElev.ID, localElev.Floor, elevatorMatrix) {
			return DIR_Up
		}
		if isOrderBelow(localElev.ID, localElev.Floor, elevatorMatrix) {
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

func shouldStop(id int, localElev Elevator, elevatorMatrix [][]int) bool {
	//Cab call is pressed, stop
	if elevatorMatrix[len(elevatorMatrix)-localElev.Floor-1][id*NumElevators+2] == 1 {
		return true
	}
	switch localElev.Dir {
	case DIR_Up:
		if elevatorMatrix[len(elevatorMatrix)-localElev.Floor-1][id*NumElevators] == 1 {
			return true
		} else if elevatorMatrix[len(elevatorMatrix)-localElev.Floor-1][id*NumElevators+1] == 1 && !isOrderAbove(id, localElev.Floor, elevatorMatrix) {
			return true
		}
	case DIR_Down:
		if elevatorMatrix[len(elevatorMatrix)-localElev.Floor-1][id*NumElevators+1] == 1 {
			return true
		} else if elevatorMatrix[len(elevatorMatrix)-localElev.Floor-1][id*NumElevators] == 1 && !isOrderBelow(id, localElev.Floor, elevatorMatrix) {
			return true
		}
	}
	return false
}
