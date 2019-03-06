package FSM

import (
	"../IO"
  "fmt"
  "../utilities"
//	"time"
  "../orderManager"

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


const numFloors = 4
const numElevators = 3


func StateMachine(FSMchans FSMchannels, elevatorMatrix [][]int){
  for {
    select {
    case newOrder := <- FSMchans.NewOrderChan:
      //Print button pressed
      fmt.Println(newOrder.Floor)
      //Add pressed button to matrix
      orderManager.AddOrder(0, elevatorMatrix, newOrder)
      // Print matrix
      utilities.PrintMatrix(elevatorMatrix,4,3)






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
        if shouldStop(0, currentFloor, elevatorMatrix) {
          elevatorMatrix[len(elevatorMatrix) - currentFloor - 1][0*3 + 2] = 0
          fmt.Println("stop")
        }
        // Print matrix
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


func updateElevFloor(elevID int, newFloor int, elevatorMatrix [][]int){
  elevatorMatrix[2][elevID*3] = newFloor
  //tror man ikke trenger å returnere matriser
}

func isOrderAbove(elevID int, currentFloor int, elevatorMatrix [][]int) bool {
  for floor := (len(elevatorMatrix) - currentFloor - 2); floor > 3; floor-- {
    for buttons := (elevID*3); buttons < (elevID*3 + 3); buttons++ {
      if elevatorMatrix[floor][buttons] == 1 {
        return true
      }
    }

  }
  return false
}

/*
func isOrderBelow(elevID int, currentFloor int, elevatorMatrix [][]int) bool {
  for floor := (len(elevatorMatrix) - currentFloor); floor > (len(elevatorMatrix) - 1); floor++ {
    for buttons := (elevID*3); buttons < (elevID*3 + 3); buttons++ {
      if elevatorMatrix[floor][buttons] == 1 {
        return true
      }
    }

  }
  return false
}
*/

func shouldStop(elevID int, currentFloor int, elevatorMatrix [][]int) bool {
  //Cab call is pressed, stop
  if elevatorMatrix[len(elevatorMatrix) - currentFloor - 1][elevID*3 + 2] == 1 {
    return true
  }

  // Also stop if elevator is going in the same direction


  return false
}











/*
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
*/
