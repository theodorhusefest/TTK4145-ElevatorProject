package syncElevator

import (
	. "../Config"
	"../Network/network/peers"
	//"../Utilities"
	"fmt"
	"strconv"
	"time"
)

type SyncElevatorChannels struct {
	OutGoingMsg     chan []Message
	InCommingMsg    chan []Message
	ChangeInOrderch chan []Message
	PeerUpdate      chan peers.PeerUpdate
	TransmitEnable  chan bool
	BroadcastTicker chan bool
}

func SyncElevator(	elevatorMatrix [][]int, elevator Elevator, 
					syncChans SyncElevatorChannels, UpdateOrderch chan Message, 
					UpdateElevStatusch chan Message, GlobalStateUpdatech chan Message, 
					MatrixUpdatech chan Message) {

	Online := false
	AckMatrix := [4 + NumFloors][3 * NumElevators]AckStruct{}
	ResendMatrixAck := AckStruct{}

	broadCastTicker := time.NewTicker(10 * time.Millisecond)
	ackTicker := time.NewTicker(500 * time.Millisecond)
	for {
		select {

		// --------------------------------------------------------------------------Case triggered by local ordermanager, change in order
		case changeInOrder := <-syncChans.ChangeInOrderch:


			if Online {


				for _, message := range changeInOrder {
					if !message.Ack {
						for elev := 0; elev < NumElevators; elev++ {
							switch message.Select {
							case NewOrder:
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].Data = 1
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] = true
							case OrderComplete:
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].Data = 0
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] = true
							case UpdateStates:

							case UpdatedMatrix:
								ResendMatrixAck.AwaitingAck[elev] = true
							}
						}
					}
				}
				select {
				case <-broadCastTicker.C:
					syncChans.OutGoingMsg <- changeInOrder
				}

			} else {
				for _, message := range changeInOrder {
					if !(message.Done) {
						message.Done = true
						switch message.Select {

						case NewOrder:
							UpdateOrderch <- message

						case OrderComplete:
							UpdateOrderch <- message

						case UpdateStates:
							UpdateElevStatusch <- message

						case UpdateOffline:
							UpdateElevStatusch <- message

						case ACK:

						case SendMatrix:
							MatrixUpdatech <- message

						case UpdatedMatrix:
							MatrixUpdatech <- message
						}
					}
				}
			}


		// --------------------------------------------------------------------------Case triggered by bcast.Recieving
		case msgRecieved := <-syncChans.InCommingMsg:

			for _, message := range msgRecieved {
				if !(message.Done) {

					sendAck := Message{Select: message.Select, Done: false, SenderID: elevator.ID, ID: message.ID, Floor: message.Floor, Button: message.Button,
						State: message.State, Dir: message.Dir, Ack: true, ResendMatrix: message.ResendMatrix, Matrix: message.Matrix}

					if message.SenderID != elevator.ID {

						message.Done = true

						switch message.Select {
						case NewOrder:
							if message.Ack {
								sendAck.Ack = false
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[message.SenderID] = false
								allReady := true
								for elev := 0; elev < NumElevators; elev++ {
									if AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
									}
								}

								if allReady {
									UpdateOrderch <- message
								}

							} else { 
								sendAck.Ack = true
								UpdateOrderch <- message
							}

						case OrderComplete:
							if message.Ack {
								sendAck.Ack = false
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[message.SenderID] = false
								allReady := true
								for elev := 0; elev < NumElevators; elev++ {
									if AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
									}
								}
								if allReady {
									UpdateOrderch <- message
								}
							} else {
								sendAck.Ack = true
								UpdateOrderch <- message
							}

						case UpdateStates:
							sendAck.Ack = false
							GlobalStateUpdatech <- message

						case UpdateOffline:
						case ACK:


						case SendMatrix:
							if message.Ack {
								sendAck.Ack = false
								ResendMatrixAck.AwaitingAck[message.SenderID] = false
								allReady := true
								for elev := 0; elev < NumElevators; elev++ {
									if ResendMatrixAck.AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
									}
								}

								if allReady {
									MatrixUpdatech <- message
								}
							} else { 
								sendAck.Ack = true
								MatrixUpdatech <- message
							}

						case UpdatedMatrix:
							if message.Ack {
								sendAck.Ack = false
								ResendMatrixAck.AwaitingAck[message.SenderID] = false
								allReady := true
								for elev := 0; elev < NumElevators; elev++ {
									if ResendMatrixAck.AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
									}
								}

								if allReady {
									MatrixUpdatech <- message
								}
							} else { 
								sendAck.Ack = true
								MatrixUpdatech <- message
							}
						}



					} else { 
						sendAck.Ack = false
						message.Done = true
						switch message.Select {
						case NewOrder:
							AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[message.SenderID] = false
							allReady := true
							for elev := 0; elev < NumElevators; elev++ {
								if AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
									allReady = false
								}
							}
							if allReady {
								fmt.Println(message)
								UpdateOrderch <- message
							}

						case OrderComplete:
							AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[message.SenderID] = false
							allReady := true
							for elev := 0; elev < NumElevators; elev++ {
								if AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
									allReady = false
								}
							}
							if allReady {
								UpdateOrderch <- message
							}

						case UpdateStates:
							GlobalStateUpdatech <- message

						case UpdateOffline:
						case ACK:
						case SendMatrix:

						case UpdatedMatrix:
							ResendMatrixAck.AwaitingAck[message.SenderID] = false
							// Hvis alle har sendt ack: kjør
							allReady := true
							for elev := 0; elev < NumElevators; elev++ {
								if ResendMatrixAck.AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
									allReady = false
								}
							}

							if allReady {
								MatrixUpdatech <- message
							}
						}

					}

					if sendAck.Ack {
						OutAck := []Message{{Select: sendAck.Select, Done: false, SenderID: elevator.ID, ID: sendAck.ID, Floor: sendAck.Floor, Button: sendAck.Button,
							State: sendAck.State, Dir: sendAck.Dir, Ack: true, ResendMatrix: sendAck.ResendMatrix, Matrix: sendAck.Matrix}}

						select {
						case <-broadCastTicker.C:
							syncChans.OutGoingMsg <- OutAck
						}

					}
				}
			}


		case <-ackTicker.C:

			for floor := 4; floor < 4+NumFloors; floor++ {
				for button := 0; button < 3*NumElevators; button++ {
					for elev := 0; elev < NumElevators; elev++ {
						if AckMatrix[floor][button].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 && elevator.ID != elev {
							//utilities.PrintAckMatrix(AckMatrix, NumFloors , NumElevators)
							ResendOrder := []Message{}
							if AckMatrix[floor][button].Data == 1 {
								ResendOrder = []Message{{Select: NewOrder, Done: false, SenderID: elevator.ID, ID: (button / NumElevators), Floor: 7 - floor, Button: ButtonType(button % 3)}}
							} else {
								ResendOrder = []Message{{Select: OrderComplete, Done: false, SenderID: elevator.ID, ID: (button / NumElevators), Floor: 7 - floor, Button: ButtonType(button % 3)}}
							}
							select {
							case <-broadCastTicker.C:
								syncChans.OutGoingMsg <- ResendOrder
							}
						}
					}
				}
			}
			for elev := 0; elev < NumElevators; elev++ {
				if ResendMatrixAck.AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
					message := Message{Select: SendMatrix, ID: elev}
					MatrixUpdatech <- message
				}
			}

		// --------------------------------------------------------------------------Case triggered by update in peers
		case p := <-syncChans.PeerUpdate:

			time.Sleep(500 * time.Millisecond)
			AckMatrix = [4 + NumFloors][3 * NumElevators]AckStruct{}



			if len(p.New) > 0 {
				newID, _ := strconv.Atoi(p.New) // ID of new Peer

				if newID == elevator.ID && len(p.Peers) == 1 {
					// You are alone on network (Either first or someone disappeard)
					// do nothing
					Online = true
					fmt.Println(newID, "is online")
				} else if newID == elevator.ID && len(p.Peers) > 1 {
					// Either been offline or first time online
					// Ask for matrix
					Online = true
					fmt.Println(newID, "is also online")

				} else if newID != elevator.ID && Online {
					// Already online, send matrix to new
					message := Message{Select: SendMatrix, ID: newID}
					MatrixUpdatech <- message
				}
			}

			for _, peerLost := range p.Lost {
				newID, _ := strconv.Atoi(peerLost)
				if newID != elevator.ID {
					// Someone else is offline
					fmt.Println(newID, "is offline")
					message := Message{Select: UpdateOffline, ID: newID, Done: true}
					GlobalStateUpdatech <- message
				} else {
					Online = false
					fmt.Println("I am offline")
				}
			}
		}
	}
}
