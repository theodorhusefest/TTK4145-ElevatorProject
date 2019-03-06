package main

import (
//  "fmt"
  "./Initialize"
  "./utilities"
  //"./orderManager"
  "./IO"
  "./FSM"
)




func main() {

  const numFloors = 4
  const numElevators = 3

  io.Init("localhost:15657",4)

  FSMchans := FSM.FSMchannels{
    NewOrderChan: make(chan io.ButtonEvent),
    ArrivedAtFloorChan: make(chan int),
    DoorTimeoutChan:  make(chan bool),
  }

  channelFloor := make(chan int) //channel that is used in InitElevator. Should maybe have a struct with channels?
  elevatorMatrix := initialize.InitializeMatrix(numFloors,numElevators)  // Set up matrix, add ID
  initialize.InitElevator(0,elevatorMatrix,channelFloor)  // Move elevator to nearest floor and update matrix


  go io.PollFloorSensor(FSMchans.ArrivedAtFloorChan)
  go io.PollButtons(FSMchans.NewOrderChan)

  go FSM.StateMachine(FSMchans, elevatorMatrix)


/*
  for {
    select {
    case buttonPressed := <- buttonChn:
      elevatorMatrix = orderManager.AddOrder(0, elevatorMatrix, buttonPressed)
      utilities.PrintMatrix(elevatorMatrix,numFloors,numElevators)

    case floor := <- floorChn:
      fmt.Println(floor)


    }
  }
*/
  utilities.PrintMatrix(elevatorMatrix,numFloors,numElevators)
  select{}
}
