package main

import (
	. "./Config"
	"./FSM"
	"./IO"
	"./Initialize"
	"./Network/network/bcast"
	"./Network/network/peers"
	"./elevatorSync"
	"./orderManager"
	"flag"
	"strconv"
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
		LocalOrderFinishedChan: make(chan int, 2),
		NewLocalOrderch:        make(chan Message, 2),
		UpdateOrderch:          make(chan Message, 2),
		MatrixUpdatech:         make(chan Message),
	}
	// Channels for SyncElevator
	SyncElevatorChans := syncElevator.SyncElevatorChannels{
		OutGoingMsg:    make(chan []Message, 2),
		InCommingMsg:   make(chan []Message, 2),
		SyncUpdatech:   make(chan []Message, 2),
		PeerUpdate:     make(chan peers.PeerUpdate, 2),
		TransmitEnable: make(chan bool),
	}
	var (
		ButtonPressedch     = make(chan ButtonEvent)
		LocalStateUpdatech  = make(chan Message)
		GlobalStateUpdatech = make(chan Message)
	)

	channelFloor := make(chan int)
	initialize.InitElevator(localElevator, elevatorMatrix, channelFloor)

	go io.PollFloorSensor(FSMchans.ArrivedAtFloorChan)
	go FSM.StateMachine(elevatorMatrix, localElevator, FSMchans,
		OrderManagerchans.LocalOrderFinishedChan,
		LocalStateUpdatech)

	go io.PollButtons(ButtonPressedch)
	go orderManager.OrderManager(elevatorMatrix, localElevator,
		OrderManagerchans, ButtonPressedch, FSMchans.NewLocalOrderChan,
		SyncElevatorChans.SyncUpdatech, GlobalStateUpdatech)

	go orderManager.UpdateElevStatus(elevatorMatrix, LocalStateUpdatech, SyncElevatorChans.SyncUpdatech, localElevator)

	go syncElevator.SyncElevator(elevatorMatrix, localElevator, SyncElevatorChans,
		OrderManagerchans.UpdateOrderch, LocalStateUpdatech,
		GlobalStateUpdatech, OrderManagerchans.MatrixUpdatech)

	go peers.Transmitter(15789, strconv.Itoa(localElevator.ID), SyncElevatorChans.TransmitEnable)
	go peers.Receiver(15789, SyncElevatorChans.PeerUpdate)

	go bcast.Transmitter(15790, SyncElevatorChans.OutGoingMsg)
	go bcast.Receiver(15790, SyncElevatorChans.InCommingMsg)

	select {}

}
