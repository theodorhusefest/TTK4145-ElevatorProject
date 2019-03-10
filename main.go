package main

import (
//  "fmt"
  . "./Config"
  "./Initialize"
  "./Utilities"
  "./orderManager"
  "./IO"
  "./FSM"
  "time"
)




func main() {

  io.Init("localhost:15657",4)

  FSMchans := FSM.FSMchannels{
    NewLocalOrderChan: make(chan int),
    ArrivedAtFloorChan: make(chan int),
    DoorTimeoutChan:  make(chan bool),
  }

  OrderManagerchans := orderManager.OrderManagerChannels{
    UpdateElevatorChan: make(chan Elevator),
    LocalOrderFinishedChan: make(chan int),
  }
  var (
    NewGlobalOrderChan = make(chan ButtonEvent)
  )

  channelFloor := make(chan int) //channel that is used in InitElevator. Should maybe have a struct with channels?
  elevatorMatrix := initialize.InitializeMatrix(NumFloors,NumElevators)  // Set up matrix, add ID
  initialize.InitElevator(0,elevatorMatrix,channelFloor)  // Move elevator to nearest floor and update matrix

  // FSM goroutines
  go io.PollFloorSensor(FSMchans.ArrivedAtFloorChan)
  go FSM.StateMachine(FSMchans, OrderManagerchans.LocalOrderFinishedChan, elevatorMatrix)

  // OrderManager goroutines
  go io.PollButtons(NewGlobalOrderChan)
  go orderManager.OrderManager(OrderManagerchans, NewGlobalOrderChan, FSMchans.NewLocalOrderChan, elevatorMatrix)



  time.Sleep(10*time.Second)
  utilities.PrintMatrix(elevatorMatrix, NumFloors, NumElevators)



/*
  for {
    select {
    case buttonPressed := <- buttonChn:
      elevatorMatrix = orderManager.AddOrder(0, elevatorMatrix, buttonPressed)
      utilities.PrintMatrix(elevatorMatrix,NumFloors,NumElevators)

    case floor := <- floorChn:
      fmt.Println(floor)


    }
  }
*/
  select{}
}
