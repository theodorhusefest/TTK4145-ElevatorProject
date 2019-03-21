package main

import (
//  "fmt"
  . "./Config"
  "./Initialize"
  //"./Utilities"
  "./orderManager"
  "./IO"
  "./FSM"
  "./elevatorSync"
  "./Network/network/peers"
  "./Network/network/bcast"
  "time"
  "strconv"
  "flag"
)


func main() {

  // Initialize
  // !!!!!!!!!!!!!! Skru av alle lys, initialiser matrisen til antall input
  floorInp := flag.Int("numFloors",4,"an int")
  elevInp := flag.Int("numElevators",3,"an int")
  portInp := flag.String("port","15657","a string")

  flag.Parse()



  elevatorMatrix, elevConfig := initialize.Initialize(*floorInp, *elevInp)

  io.Init("localhost:" + (*portInp), NumFloors)

  // Channels for FSM
  FSMchans := FSM.FSMchannels{
    NewLocalOrderChan: make(chan int),
    ArrivedAtFloorChan: make(chan int),
    DoorTimeoutChan:  make(chan bool),
  }
  // Channels for OrderManager
  OrderManagerchans := orderManager.OrderManagerChannels{
    LocalOrderFinishedChan: make(chan int),
    NewLocalOrderch: make(chan Message),
    UpdateOrderch: make(chan Message),
    MatrixUpdatech: make(chan Message),

  }
  // Channels for SyncElevator
  SyncElevatorChans := syncElevator.SyncElevatorChannels{
    OutGoingMsg: make(chan []Message),
    InCommingMsg: make(chan []Message),
    ChangeInOrderch: make(chan []Message),
    PeerUpdate: make(chan peers.PeerUpdate),
    TransmitEnable: make(chan bool),
    BroadcastTicker: make(chan bool),
  }
  var (
    NewGlobalOrderChan = make(chan ButtonEvent)
    UpdateElevStatusch = make(chan Message)
  )

  channelFloor := make(chan int) //channel that is used in InitElevator. Should maybe have a struct with channels?
  //elevatorMatrix := initialize.InitializeMatrix(NumFloors,NumElevators)  // Set up matrix, add ID
  initialize.InitElevator(elevConfig,elevatorMatrix,channelFloor)  // Move elevator to nearest floor and update matrix
  //utilities.PrintMatrix(elevatorMatrix, elevConfig.NumFloors,elevConfig.NumElevators)

  // Goroutines used in FSM
  go io.PollFloorSensor(FSMchans.ArrivedAtFloorChan)
  go FSM.StateMachine(FSMchans, OrderManagerchans.LocalOrderFinishedChan, UpdateElevStatusch, elevatorMatrix, elevConfig)

  // Goroutines used in OrderManager
  go io.PollButtons(NewGlobalOrderChan)
  go orderManager.OrderManager(elevatorMatrix, elevConfig, OrderManagerchans, NewGlobalOrderChan, FSMchans.NewLocalOrderChan, SyncElevatorChans.OutGoingMsg,
        SyncElevatorChans.ChangeInOrderch, UpdateElevStatusch)

  go orderManager.UpdateElevStatus(elevatorMatrix, UpdateElevStatusch, SyncElevatorChans.ChangeInOrderch)

  // Goroutines used in SyncElevator

  go syncElevator.SyncElevator(elevatorMatrix,SyncElevatorChans, elevConfig, OrderManagerchans.UpdateOrderch, UpdateElevStatusch, OrderManagerchans.MatrixUpdatech)

  // Goroutines used in Network/Peers
  go peers.Transmitter(15789, strconv.Itoa(elevConfig.ElevID), SyncElevatorChans.TransmitEnable)
  go peers.Receiver(15789, SyncElevatorChans.PeerUpdate)

  //  Goroutines used in Network/Bcast
  go bcast.Transmitter(15790, SyncElevatorChans.OutGoingMsg)
  go bcast.Receiver(15790, SyncElevatorChans.InCommingMsg)


  time.Sleep(10*time.Second)

  select{}


}


