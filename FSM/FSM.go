package FSM

import (
	. "../Config"
	"../IO"
	"../orderManager"
	"fmt"
	"time"
)

type FSMchannels struct {
	NewLocalOrderChan  chan int
	ArrivedAtFloorChan chan int
	DoorTimeoutChan    chan bool
}

func StateMachine(elevatorMatrix [][]int, localElev Elevator,
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

		case <-motorFailureTimeOut.C:
			fmt.Println("Motor has failed")
			localElev.State = UNDEFINED
			orderManager.InsertState(localElev.ID, int(UNDEFINED), elevatorMatrix)

			updatedStates := Message{Select: UpdateStates, ID: localElev.ID, State: int(localElev.State), Floor: localElev.Floor, Dir: localElev.Dir}
			UpdateElevStatusch <- updatedStates
		}
	}
}
