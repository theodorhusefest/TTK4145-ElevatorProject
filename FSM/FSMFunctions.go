package FSM

import (
	. "../Config"
)

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
