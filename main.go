package main

import (
//"fmt"
  . "./Config"
  "./Initialize"
  "./Utilities"
  "./orderManager"
  "./IO"
  "./FSM"
  "./elevatorSync"
  "./Network/network/peers"
  "./Network/network/bcast"
  "time"
  "strconv"

)




func main() {

  elevatorMatrix, elevConfig := initialize.Initialize()




  io.Init("localhost:15657",4)



  FSMchans := FSM.FSMchannels{
    NewLocalOrderChan: make(chan int),
    ArrivedAtFloorChan: make(chan int),
    DoorTimeoutChan:  make(chan bool),
  }
  OrderManagerchans := orderManager.OrderManagerChannels{
    UpdateElevatorChan: make(chan Message),
    LocalOrderFinishedChan: make(chan int),
  }
  SyncElevatorChans := syncElevator.SyncElevatorChannels{
    OutGoingMsg: make(chan Message),
    InCommingMsg: make(chan Message),
    ChangeInOrderch: make(chan Message),
    PeerUpdate: make(chan peers.PeerUpdate),
    TransmitEnable: make(chan bool),
    BroadcastTicker: make(chan bool),
  }
  var (
    NewGlobalOrderChan = make(chan ButtonEvent)
  )

  channelFloor := make(chan int) //channel that is used in InitElevator. Should maybe have a struct with channels?
  //elevatorMatrix := initialize.InitializeMatrix(NumFloors,NumElevators)  // Set up matrix, add ID
  initialize.InitElevator(elevConfig,elevatorMatrix,channelFloor)  // Move elevator to nearest floor and update matrix
  utilities.PrintMatrix(elevatorMatrix, elevConfig.NumFloors,elevConfig.NumElevators)





  // FSM goroutines
  go io.PollFloorSensor(FSMchans.ArrivedAtFloorChan)
  go FSM.StateMachine(FSMchans, OrderManagerchans.LocalOrderFinishedChan, elevatorMatrix,elevConfig)

  // OrderManager goroutines
  go io.PollButtons(NewGlobalOrderChan)
  go orderManager.OrderManager(OrderManagerchans, NewGlobalOrderChan, FSMchans.NewLocalOrderChan, elevatorMatrix, SyncElevatorChans.OutGoingMsg, SyncElevatorChans.ChangeInOrderch, elevConfig)


  //Sync
  go syncElevator.SyncElevator(SyncElevatorChans, elevConfig, OrderManagerchans.UpdateElevatorChan)

  //Update peers
  go peers.Transmitter(15789, strconv.Itoa(elevConfig.ElevID), SyncElevatorChans.TransmitEnable)
  go peers.Receiver(15789, SyncElevatorChans.PeerUpdate)

  //Send/recieve orders
  go bcast.Transmitter(15790, SyncElevatorChans.OutGoingMsg)
  go bcast.Receiver(15790, SyncElevatorChans.InCommingMsg)







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
