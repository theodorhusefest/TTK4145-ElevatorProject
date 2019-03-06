package FSM

import (
	"../IO"
  "fmt"
  "../utilities"
//	"time"
//  "../orderManager"

)

type state int

const (
  IDLE state = 0
  MOVING state = 1
  DOOROPEN state = 2
)

type FSMchannels struct{
  NewOrderChan chan io.ButtonEvent
  ArrivedAtFloorChan chan int
  DoorTimeoutChan chan bool
} //flytt til egen fil


func StateMachine(FSMchans FSMchannels, elevatorMatrix [][]int){
  for {
    select {
    case newOrder := <- FSMchans.NewOrderChan:
      fmt.Println(newOrder.Floor)
      /*
        If IDLE -> move to floor
          Choose direction and go
          */
          /*if (isOrderAbove()){
          } else if (isOrderBelow()){
            io.SetMotorDirection(io.MS_Down)
          }*/
          //else: Do nothing?

      /*
        If Moving
          Check if order is on the way, stop if true.
        if Opendoor -> do nothing

      */

    case currentFloor := <- FSMchans.ArrivedAtFloorChan:

        //Update floor in matrix
        updateElevFloor(0,currentFloor,elevatorMatrix)
        utilities.PrintMatrix(elevatorMatrix,4,3)
      /*
        If shouldStop -> stop elevator
          opendoorLamp
          clear order in matrix
          timeout <- doorTimeoutChan
      */


    //case  := <- FSMchans.doorTimeoutChan:
      /*
        closedoorLamp
        state = IDLE
      */
    }
  }
}


func updateElevFloor(elevID int, newFloor int, elevMatrix [][]int){
  elevMatrix[2][elevID*3] = newFloor
  //tror man ikke trenger å returnere matriser
}


func isOrderAbove(elevID int, order int, elevMatrix [][]int) bool{
  if (order > elevMatrix[elevID*3][2]){
    return true
  }
  return false
}

func isOrderBelow(elevID int, order int, elevMatrix [][]int) bool{
  if (order < elevMatrix[elevID*3][2]){
    return true
  }
  return false
}
