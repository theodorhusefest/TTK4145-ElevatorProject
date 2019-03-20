package FSM

import (
	"../IO"
  "fmt"
  "../orderManager"
	. "../Config"
  "time"
  "../Utilities"

)


type FSMchannels struct{
  NewLocalOrderChan chan int
  ArrivedAtFloorChan chan int
  DoorTimeoutChan chan bool
}

func StateMachine(FSMchans FSMchannels, LocalOrderFinishedChan chan int, UpdateElevStatusch chan []Message, elevatorMatrix [][]int, elevatorConfig ElevConfig){
	elevator := Elevator{
			ID: elevatorMatrix[0][elevatorConfig.ElevID*3],
			State: IDLE,
			Floor: elevatorMatrix[2][elevatorConfig.ElevID*3],
			Dir: DIR_Stop,
	}

	doorOpenTimeOut := time.NewTimer(3 * time.Second)
	doorOpenTimeOut.Stop()


  for {
    select {
    case newLocalOrder := <- FSMchans.NewLocalOrderChan:
			fmt.Println("he")

			switch elevator.State {
			case IDLE:
				//fmt.Println(IDLE)
				orderManager.InsertState(elevatorConfig.ElevID, int(IDLE), elevatorMatrix)
				elevator.Dir = chooseDirection(elevatorConfig.ElevID, elevatorMatrix, elevator)
				io.SetMotorDirection(elevator.Dir)
				orderManager.InsertDirection(elevatorConfig.ElevID, elevator.Dir, elevatorMatrix)
				if elevator.Dir == DIR_Stop {
					// Open door for 3 seconds
					io.SetDoorOpenLamp(true)
					//clearFloors(elevator,elevatorMatrix)
					elevator.State = DOOROPEN
					orderManager.InsertState(elevatorConfig.ElevID, int(DOOROPEN), elevatorMatrix)
					doorOpenTimeOut.Reset(3 * time.Second)
					LocalOrderFinishedChan <- elevator.Floor
				}
				elevator.State = MOVING
				orderManager.InsertState(elevatorConfig.ElevID, int(MOVING), elevatorMatrix)


			case MOVING:


			case DOOROPEN:
				if elevator.Floor == newLocalOrder {
					doorOpenTimeOut.Reset(3 * time.Second)
					fmt.Println("Resetting Time")
					LocalOrderFinishedChan <- elevator.Floor
				}
			}

			fmt.Println("End of IDLE")

			updatedElev := []Message{{Select: 3, ID: elevator.ID ,State: int(elevator.State) ,Floor: elevator.Floor, Dir: elevator.Dir}}
			UpdateElevStatusch <- updatedElev




    case currentFloor := <- FSMchans.ArrivedAtFloorChan:
        //Update floor in matrix
        orderManager.InsertFloor(elevatorConfig.ElevID,currentFloor,elevatorMatrix)
		elevator.Floor = currentFloor
		io.SetFloorIndicator(currentFloor)

		utilities.PrintMatrix(elevatorMatrix, elevatorConfig.NumFloors,elevatorConfig.NumElevators)

        if shouldStop(elevatorConfig.ElevID, elevator, elevatorMatrix) {
			elevator.State = DOOROPEN
			orderManager.InsertState(elevatorConfig.ElevID, int(DOOROPEN), elevatorMatrix)
			io.SetDoorOpenLamp(true)
			doorOpenTimeOut.Reset(3 * time.Second)

			io.SetMotorDirection(DIR_Stop)
			orderManager.InsertDirection(elevatorConfig.ElevID, elevator.Dir, elevatorMatrix)

			LocalOrderFinishedChan <- elevator.Floor
        }
		updatedElev := []Message{{Select: 3, ID: elevator.ID ,State: int(elevator.State) ,Floor: elevator.Floor, Dir: elevator.Dir}}
		UpdateElevStatusch <- updatedElev


		case <-doorOpenTimeOut.C:
			fmt.Println("DOOROPEN")
			io.SetDoorOpenLamp(false)
			elevator.Dir = chooseDirection(elevatorConfig.ElevID,elevatorMatrix, elevator)
			orderManager.InsertDirection(elevatorConfig.ElevID, elevator.Dir, elevatorMatrix)
			io.SetMotorDirection(elevator.Dir)
			LocalOrderFinishedChan <- elevator.Floor
			if elevator.Dir == DIR_Stop {
				elevator.State = IDLE
				orderManager.InsertState(elevatorConfig.ElevID, int(IDLE), elevatorMatrix)
			} else {
				//io.SetMotorDirection(elevator.Dir)
				elevator.State = MOVING
				orderManager.InsertState(elevatorConfig.ElevID, int(MOVING), elevatorMatrix)
				fmt.Println(MOVING)
			}
			updatedElev := []Message{{Select: 3, ID: elevator.ID ,State: int(elevator.State) ,Floor: elevator.Floor, Dir: elevator.Dir}}
			UpdateElevStatusch <- updatedElev
    }

  }
}

func matrixIsEmpty(elevatorMatrix [][]int, elevatorConfig ElevConfig) bool{
	for floor := 4; floor < 4+ NumFloors; floor++ {
    for buttons := (elevatorConfig.ElevID*3); buttons < (elevatorConfig.ElevID*3 + 3); buttons++ {
      if elevatorMatrix[floor][buttons] == 1 {
        return false
      }
    }
  }
	return true
}

func isOrderAbove(ElevID int, currentFloor int, elevatorMatrix [][]int) bool {
	if currentFloor == 3 {
		return false
	}

	for floor := (len(elevatorMatrix) - currentFloor - 2); floor > 3; floor-- {
    for buttons := (ElevID*3); buttons < (ElevID*3 + 3); buttons++ {
      if elevatorMatrix[floor][buttons] == 1 {
        return true
      }
    }
  }
  return false
}


func isOrderBelow(ElevID int, currentFloor int, elevatorMatrix [][]int) bool {
	if currentFloor == 0 {
		return false
	}

	for floor := (len(elevatorMatrix) - currentFloor); floor < (len(elevatorMatrix)); floor++ {
    for buttons := (ElevID*3); buttons < (ElevID*3 + 3); buttons++ {
      if elevatorMatrix[floor][buttons] == 1 {
        return true
      }
    }

  }
  return false
}

func shouldStop(id int, elevator Elevator, elevatorMatrix [][]int) bool {
  //Cab call is pressed, stop
  if elevatorMatrix[len(elevatorMatrix) - elevator.Floor - 1][id*NumElevators + 2] == 1 {
    return true
  }
	switch elevator.Dir{
	case DIR_Up:
		if elevatorMatrix[len(elevatorMatrix) - elevator.Floor - 1][id*NumElevators] == 1 {
	    return true
	  } else if elevatorMatrix[len(elevatorMatrix) - elevator.Floor - 1][id*NumElevators + 1] == 1 && !isOrderAbove(id, elevator.Floor, elevatorMatrix) {
			return true
		}
	case DIR_Down:
		if elevatorMatrix[len(elevatorMatrix) - elevator.Floor - 1][id*NumElevators + 1] == 1 {
	    return true
	  } else if elevatorMatrix[len(elevatorMatrix) - elevator.Floor - 1][id*NumElevators] == 1 && !isOrderBelow(id, elevator.Floor, elevatorMatrix) {
			return true
		}
	}
  return false
}

func chooseDirection(id int,elevatorMatrix [][]int, elevator Elevator) MotorDirection {

	if isOrderAbove(id, elevator.Floor, elevatorMatrix){
		return DIR_Up
	}
	if isOrderBelow(id, elevator.Floor, elevatorMatrix){
		return DIR_Down
	}
	return DIR_Stop
}
