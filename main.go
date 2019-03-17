package main

import (
//  "fmt"
  . "./Config"
  "./Initialize"
  "./Utilities"
  "./orderManager"
  "./IO"
  "./FSM"
  "./elevatorSync"
  "./Network/network/peers"
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
  SyncElevatorChans := syncElevator.SyncElevatorChannels{
  //  OutGoingOrder: make(chan ??)
    //InComingOrder: make(chan ??)
    PeerUpdate: make(chan peers.PeerUpdate),
    TransmitEnable: make(chan bool),
    BroadcastTicker: make(chan bool),
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


  //Sync
  go syncElevator.SyncElevator(SyncElevatorChans)
  go peers.Transmitter(15789, string("Heis0"), SyncElevatorChans.TransmitEnable)
  go peers.Transmitter(15789, string("Heis1"), SyncElevatorChans.TransmitEnable)
  go peers.Receiver(15789, SyncElevatorChans.PeerUpdate)



  time.Sleep(10*time.Second)
  utilities.PrintMatrix(elevatorMatrix, NumFloors, NumElevators)

  //initialize network module

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
