package main

import (
	//  "fmt"
	. "./Config"
	"./Initialize"
	//"./Utilities"
	"./FSM"
	"./IO"
	"./Network/network/bcast"
	"./Network/network/peers"
	"./elevatorSync"
	"./orderManager"
	"flag"
	"strconv"
	"time"
)

func main() {
	
	floorInp := flag.Int("numFloors", 4, "an int")
	elevInp := flag.Int("numElevators", 3, "an int")
	portInp := flag.String("port", "15657", "a string")

	flag.Parse()

	elevatorMatrix, localElevator := initialize.Initialize(*floorInp, *elevInp)


	io.Init("localhost:"+(*portInp), NumFloors)

	// Channels for FSM
	FSMchans := FSM.FSMchannels{
		NewLocalOrderChan:  make(chan int),
		ArrivedAtFloorChan: make(chan int),
		DoorTimeoutChan:    make(chan bool),
	}
	// Channels for OrderManager
	OrderManagerchans := orderManager.OrderManagerChannels{
		LocalOrderFinishedChan: make(chan int,2),
		NewLocalOrderch:        make(chan Message, 2),
		UpdateOrderch:          make(chan Message, 2),
		MatrixUpdatech:         make(chan Message),
	}
	// Channels for SyncElevator
	SyncElevatorChans := syncElevator.SyncElevatorChannels{
		OutGoingMsg:     make(chan []Message),
		InCommingMsg:    make(chan []Message),
		ChangeInOrderch: make(chan []Message, 2),
		PeerUpdate:      make(chan peers.PeerUpdate),
		TransmitEnable:  make(chan bool),
		BroadcastTicker: make(chan bool),
	}
	var (
		ButtonPressedch = make(chan ButtonEvent)
		UpdateElevStatusch = make(chan Message)
		GlobalStateUpdatech = make(chan Message)
	)

	channelFloor := make(chan int)
	initialize.InitElevator(localElevator, elevatorMatrix, channelFloor)

	// Goroutines used in FSM
	go io.PollFloorSensor(FSMchans.ArrivedAtFloorChan)
	go FSM.StateMachine(FSMchans, OrderManagerchans.LocalOrderFinishedChan, UpdateElevStatusch, elevatorMatrix, localElevator)

	// Goroutines used in OrderManager
	go io.PollButtons(ButtonPressedch)
	go orderManager.OrderManager(elevatorMatrix, localElevator, OrderManagerchans, ButtonPressedch, FSMchans.NewLocalOrderChan,
		SyncElevatorChans.ChangeInOrderch, UpdateElevStatusch, GlobalStateUpdatech)

	go orderManager.UpdateElevStatus(elevatorMatrix, UpdateElevStatusch, SyncElevatorChans.ChangeInOrderch, localElevator)

	// Goroutines used in SyncElevator
	go syncElevator.SyncElevator(elevatorMatrix, localElevator, SyncElevatorChans, OrderManagerchans.UpdateOrderch, UpdateElevStatusch, GlobalStateUpdatech, OrderManagerchans.MatrixUpdatech)

	// Goroutines used in Network/Peers
	go peers.Transmitter(15789, strconv.Itoa(localElevator.ID), SyncElevatorChans.TransmitEnable)
	go peers.Receiver(15789, SyncElevatorChans.PeerUpdate)

	//  Goroutines used in Network/Bcast
	go bcast.Transmitter(15790, SyncElevatorChans.OutGoingMsg)
	go bcast.Receiver(15790, SyncElevatorChans.InCommingMsg)

	time.Sleep(10 * time.Second)

	select {}

}
