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

func SyncElevator(elevatorMatrix [][]int, syncChans SyncElevatorChannels, elevator Elevator,
	UpdateOrderch chan Message, UpdateElevStatusch chan Message, UpdateOfflinech chan Message, MatrixUpdatech chan Message) {

	Online := false
	AckMatrix := [4 + NumFloors][3 * NumElevators]AckStruct{}
	ResendMatrixAck := AckStruct{}
	//utilities.PrintAckMatrix(AckMatrix, NumFloors , NumElevators)

	broadCastTicker := time.NewTicker(10 * time.Millisecond)
	ackTicker := time.NewTicker(1500 * time.Millisecond)
	for {
		select {

		// --------------------------------------------------------------------------Case triggered by local ordermanager, change in order
		case changeInOrder := <-syncChans.ChangeInOrderch:
			//Håndter endring som kom fra ordermanager. Send alt inn på message og sett message.Done = false
			//message := changeInOrder
			//utilities.PrintMatrix(elevatorMatrix, 4,3)

			if Online {

				// If outgoing msg is an ACK msg: Do nothing
				// else: You are sending something new, add new order to matrix and set all awaiting to true, recieving to false
				// Send out order

				for _, message := range changeInOrder {
					if !message.Ack {
						fmt.Println("Sending message:", MessageType(message.Select), message.ID, message.SenderID)
						for elev := 0; elev < NumElevators; elev++ {
							switch message.Select {
							case NewOrder:
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].Data = 1
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] = true
							case OrderComplete:
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].Data = 0
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] = true
							case UpdateStates:
								//AckMatrix[1][message.SenderID*3].Data = message.State
								//AckMatrix[1][message.SenderID*3].AwaitingAck[elev] = true
							case UpdatedMatrix:
								ResendMatrixAck.AwaitingAck[elev] = true
							}
						}
					}
				}
				select {
				case <-broadCastTicker.C:
					//fmt.Println(elevator.ID, "is sending outgoing message ")
					syncChans.OutGoingMsg <- changeInOrder
				}

			} else {
				for _, message := range changeInOrder {
					if !(message.Done) {
						message.Done = true
						switch message.Select {

						case NewOrder:
							UpdateOrderch <- message
							fmt.Println("UpdateOrderch 4")

						case OrderComplete:
							UpdateOrderch <- message

						case UpdateStates:
							UpdateElevStatusch <- message

						case UpdateOffline:
							UpdateElevStatusch <- message

						case ACK:
							// Not online, cannot send ack

						case SendMatrix:
							MatrixUpdatech <- message

						case UpdatedMatrix:
							MatrixUpdatech <- message
						}
					}
				}
			}

			// Broadcast message

			// Vent til alle er enige, gi klarsignal til ordermanager ??????

		// --------------------------------------------------------------------------Case triggered by bcast.Recieving
		case msgRecieved := <-syncChans.InCommingMsg:

			for _, message := range msgRecieved {
				if !(message.Done) {

					// Lag en lokal struct
					sendAck := Message{Select: message.Select, Done: false, SenderID: elevator.ID, ID: message.ID, Floor: message.Floor, Button: message.Button,
						State: message.State, Dir: message.Dir, Ack: true, ResendMatrix: message.ResendMatrix, Matrix: message.Matrix}

					// Du har ikke sendt mld som kom inn. Send Ack tilbake og execute
					if message.SenderID != elevator.ID {

						message.Done = true

						switch message.Select {
						case NewOrder:
							if message.Ack {
								sendAck.Ack = false
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[message.SenderID] = false
								// Hvis alle har sendt ack: kjør
								allReady := true
								for elev := 0; elev < NumElevators; elev++ {
									if AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
									}
								}

								if allReady {
									//fmt.Println("Someone else, all acked", message)
									fmt.Println("UpdateOrderch 1")
									fmt.Println(message)
									UpdateOrderch <- message
								}

							} else { //Send ack tilbake og execute
								sendAck.Ack = true
								//fmt.Println("Someone else, not an ack", message)
								fmt.Println("UpdateOrderch 2")
								UpdateOrderch <- message
							}

						// Samme som NewOrder
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
							fmt.Println("----------")
							sendAck.Ack = false
							UpdateOfflinech <- message
							/*	if message.Ack {
									sendAck.Ack = false
									AckMatrix[1][message.ID*3].AwaitingAck[message.SenderID] = false
									// Hvis alle har sendt ack: kjør
									allReady := true
									for elev := 0; elev < NumElevators; elev++ {
										if AckMatrix[len(AckMatrix) - message.Floor - 1][message.ID*3 + int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
											allReady = false
										}
									}

									if allReady {
										//fmt.Println("Someone else, all acked", message)
										UpdateElevStatusch <- message
									}
								} else { //Send ack tilbake og execute
									sendAck.Ack = true
									//fmt.Println("Someone else, not an ack", message)
									UpdateElevStatusch <- message
								}*/
						case UpdateOffline:
						case ACK:

						// Du får inn sendmatrix som ikke er deg selv
						case SendMatrix:
							if message.Ack {
								sendAck.Ack = false
								ResendMatrixAck.AwaitingAck[message.SenderID] = false
								// Hvis alle har sendt ack: kjør
								allReady := true
								for elev := 0; elev < NumElevators; elev++ {
									if ResendMatrixAck.AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
									}
								}

								if allReady {
									//fmt.Println("Someone else, all acked", message)
									MatrixUpdatech <- message
								}
							} else { //Send ack tilbake og execute
								sendAck.Ack = true
								//fmt.Println("Someone else, not an ack", message)
								MatrixUpdatech <- message
							}

						case UpdatedMatrix:
							if message.Ack {
								sendAck.Ack = false
								ResendMatrixAck.AwaitingAck[message.SenderID] = false
								// Hvis alle har sendt ack: kjør
								allReady := true
								for elev := 0; elev < NumElevators; elev++ {
									if ResendMatrixAck.AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
									}
								}

								if allReady {
									//fmt.Println("Someone else, all acked", message)
									MatrixUpdatech <- message
								}
							} else { //Send ack tilbake og execute
								sendAck.Ack = true
								//fmt.Println("Someone else, not an ack", message)
								MatrixUpdatech <- message
							}
						}

					} else { // Du har sendt meldingen som kom inn. Legg til i AckMatrix. Sjekk om du har fått ack fra alle, isåfall utfør
						sendAck.Ack = false
						message.Done = true
						switch message.Select {
						case NewOrder:
							// Add to matrix, check if everyone has given ack. If so, run, if not, do nothing.
							AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[message.SenderID] = false
							allReady := true
							for elev := 0; elev < NumElevators; elev++ {
								if AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
									allReady = false
								}
							}
							if allReady {
								//fmt.Println("Myself, all acked ", message)
								fmt.Println("UpdateOrderch 3")
								fmt.Println(message)
								UpdateOrderch <- message
							}

						case OrderComplete:
							// Add to matrix, check if everyone has given ack. If so, run, if not, do nothing.
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
							fmt.Println("--")
							UpdateOfflinech <- message
							/*AckMatrix[1][message.ID*3].AwaitingAck[message.SenderID] = false
							// Hvis alle har sendt ack: kjør
							allReady := true
							for elev := 0; elev < NumElevators; elev++ {
								if AckMatrix[len(AckMatrix) - message.Floor - 1][message.ID*3 + int(message.Button)].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 {
										allReady = false
								}
							}

							if allReady {
									//fmt.Println("Someone else, all acked", message)
								UpdateElevStatusch <- message
							}*/

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
								//fmt.Println("Someone else, all acked", message)
								MatrixUpdatech <- message
							}
						}

					}

					if sendAck.Ack {
						OutAck := []Message{{Select: sendAck.Select, Done: false, SenderID: elevator.ID, ID: sendAck.ID, Floor: sendAck.Floor, Button: sendAck.Button,
							State: sendAck.State, Dir: sendAck.Dir, Ack: true, ResendMatrix: sendAck.ResendMatrix, Matrix: sendAck.Matrix}}

						// Send ACK
						//fmt.Println("Outgoing ack", OutAck)
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
							fmt.Println("ResendOrder: ", ResendOrder, " Button/NumElevators", button/NumElevators, button, NumElevators)
							fmt.Println(AckMatrix[floor][button].AwaitingAck)
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
					fmt.Println(ResendMatrixAck)
					message := Message{Select: SendMatrix, ID: elev}
					MatrixUpdatech <- message
				}
			}

			//fmt.Println(AckMatrix)
			//utilities.PrintAckMatrix(AckMatrix, NumFloors , NumElevators)

		// --------------------------------------------------------------------------Case triggered by update in peers
		case p := <-syncChans.PeerUpdate:

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
					UpdateOfflinech <- message
				} else {
					Online = false
					fmt.Println("I am offline")
				}
			}
		}
	}
}

/*

func addAck (matrix [4+NumFloors][3*NumElevators]AckStruct, resendMatrixAck AckStruct, messages []Message, data int) {
	for _, message := range messages {

		for elev := 0; elev < NumElevators; elev++ {
			switch message.Select {

			case NewOrder:
				matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].Data = data
				matrix[7][0].AwaitingAck[elev] = true
				matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].RecievedAck[elev] = false

			case OrderComplete:
				for i := 0; i < 3; i++ {
					matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + i].Data = 0
					matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + i].AwaitingAck[elev] = true
					matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + i].RecievedAck[elev] = false
				}

			case UpdateStates:
				matrix[1][message.SenderID*3].Data = message.State
				matrix[1][message.SenderID*3].AwaitingAck[elev] = true
				matrix[1][message.SenderID*3].RecievedAck[elev] = false

			case UpdateOffline:


			case ACK:
				//Just send ack

			case SendMatrix:



			case UpdatedMatrix:
				resendMatrixAck.AwaitingAck[elev] = true
				resendMatrixAck.RecievedAck[elev] = false
			}
		}
		fmt.Println(matrix[7][0])
	}
}*/

func addSingleAck(matrix [4 + NumFloors][3 * NumElevators]AckStruct, resendMatrixAck AckStruct, message Message) {
	switch message.Select {
	case NewOrder:
		matrix[0][0].Data = 1

	case OrderComplete:
		matrix[len(matrix)-message.Floor-1][message.SenderID*3+int(message.Button)].AwaitingAck[message.SenderID] = false

	case UpdateStates:
		matrix[1][message.SenderID*3].AwaitingAck[message.SenderID] = false

	//case UpdateOffline:
	//case ACK:
	//case SendMatrix:

	case UpdatedMatrix:
		resendMatrixAck.AwaitingAck[message.SenderID] = false
	}
}

// Returnerer false dersom det ikke ventes på ack
func checkAck(matrix [4 + NumFloors][3 * NumElevators]AckStruct, resendMatrixAck AckStruct, message Message) bool {
	awaiting := false
	for elev := 0; elev < NumElevators; elev++ {
		switch message.Select {
		case NewOrder:
			if matrix[len(matrix)-message.Floor-1][message.SenderID*3+int(message.Button)].AwaitingAck[elev] == true {
				awaiting = true
			}

		case OrderComplete:
			if matrix[len(matrix)-message.Floor-1][message.SenderID*3+int(message.Button)].AwaitingAck[elev] == true {
				awaiting = true
			}

		case UpdateStates:
			if matrix[1][message.SenderID*3].AwaitingAck[elev] == true {
				awaiting = true
			}

		//case UpdateOffline:
		//case ACK:
		//case SendMatrix:

		case UpdatedMatrix:
			if resendMatrixAck.AwaitingAck[elev] == true {
				awaiting = true
			}

		}
	}
	return awaiting
}
