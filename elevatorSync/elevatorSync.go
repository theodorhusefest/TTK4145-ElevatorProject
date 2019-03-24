package syncElevator

import (
	. "../Config"
	"../Network/network/peers"
	"fmt"
	"strconv"
	"time"
  //"../Utilities"
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
	UpdateOrderch chan Message, UpdateElevStatusch chan Message, MatrixUpdatech chan Message) {

	Online := false
	AckMatrix := [4+NumFloors][3*NumElevators]AckStruct{}
	ResendMatrixAck := AckStruct{}


	broadCastTicker := time.NewTicker(100 * time.Millisecond)
	ackTicker := time.NewTicker(20000 * time.Millisecond)
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
      					addAck(AckMatrix, ResendMatrixAck, changeInOrder, 1)
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
					sendAck := Message{Select: message.Select, Done: false, SenderID: elevator.ID, ID: message.ID, Floor: message.Floor, Button: message.Button, 
						State: message.State, Dir: message.Dir, Ack: message.Ack, ResendMatrix: message.ResendMatrix, Matrix: message.Matrix}
					if message.SenderID != elevator.ID {

						// Mottatt melding fra noen andre
						// if: Ack == true
							// --> Legg til i matrise
							// Hvis alle er lagt til, utfør ordren. Sett awaiting hos alle til false
								// Hvis ikke, vent til alle ha sendt ack evt resend ordre
						// else: ny ordre
							// Utfør ordre
							// Lag en ackMsg og send tilbake

						message.Done = true

						switch message.Select {
							case NewOrder:
								if message.Ack {
									sendAck.Ack = false
									addSingleAck(AckMatrix, ResendMatrixAck, message)								
									if !checkAck(AckMatrix, ResendMatrixAck, message) {
										UpdateOrderch <- message
									}
								} else {
								sendAck.Ack = true
								UpdateOrderch <- message
								}

							case OrderComplete:
								if message.Ack {
									sendAck.Ack = false
									addSingleAck(AckMatrix, ResendMatrixAck, message)
									if !checkAck(AckMatrix, ResendMatrixAck, message) {
										UpdateOrderch <- message
									}
								} else {
								sendAck.Ack = true
								UpdateOrderch <- message
								}

							case UpdateStates:
								if message.Ack {
									sendAck.Ack = false
									addSingleAck(AckMatrix, ResendMatrixAck, message)
									if !checkAck(AckMatrix, ResendMatrixAck, message) {
										UpdateElevStatusch <- message
									}
								} else {
								sendAck.Ack = true
								UpdateElevStatusch <- message
								}

							case UpdateOffline:
							case ACK:
							case SendMatrix:

							case UpdatedMatrix:
								if message.Ack {
									sendAck.Ack = false
									addSingleAck(AckMatrix, ResendMatrixAck, message)
									if !checkAck(AckMatrix, ResendMatrixAck, message) {
										MatrixUpdatech <- message
									}
								} else {
								sendAck.Ack = true
								MatrixUpdatech <- message
								}


						} 
					} else {
							// Mottatt melding fra deg selv
							// if: Ack == true
								// --> Legg til i matrise
								// Hvis alle er lagt til, utfør ordren
						// Skal ikke send noen ack tilbake
						sendAck.Ack = false

						switch message.Select {
						case NewOrder:
							if message.Ack {
								addSingleAck(AckMatrix, ResendMatrixAck, message)
								if !checkAck(AckMatrix, ResendMatrixAck, message) {
									UpdateOrderch <- message
								}
							} 

						case OrderComplete:
							if message.Ack {
								addSingleAck(AckMatrix, ResendMatrixAck, message)
								if !checkAck(AckMatrix, ResendMatrixAck, message) {
									UpdateOrderch <- message									
								}
							} 

						case UpdateStates:
							if message.Ack {
								addSingleAck(AckMatrix, ResendMatrixAck, message)
								if !checkAck(AckMatrix, ResendMatrixAck, message) {
									UpdateElevStatusch <- message
								}
							} 

						case UpdateOffline:
						case ACK:
						case SendMatrix:
						
						case UpdatedMatrix:
							if message.Ack {
								addSingleAck(AckMatrix, ResendMatrixAck, message)
								if !checkAck(AckMatrix, ResendMatrixAck, message) {
									MatrixUpdatech <- message
								}
							} 
						} 
					}
					if sendAck.Ack {
						OutAck := []Message{{Select: sendAck.Select, Done: false, SenderID: elevator.ID, ID: sendAck.ID, Floor: sendAck.Floor, Button: sendAck.Button, 
							State: sendAck.State, Dir: sendAck.Dir, Ack: sendAck.Ack, ResendMatrix: sendAck.ResendMatrix, Matrix: sendAck.Matrix}}
						syncChans.ChangeInOrderch <- OutAck
					}
				}
			}

















			/*WORKING WITHOUT ACK

			for _, message := range msgRecieved {
				if !(message.Done) {
					message.Done = true
					switch message.Select {

					case NewOrder:
						UpdateOrderch <- message
						message.Ack = true

					case OrderComplete:
						UpdateOrderch <- message
						message.Ack = true

					case UpdateStates:
						UpdateElevStatusch <- message
						message.Ack = true

					case UpdateOffline:
						UpdateElevStatusch <- message
						message.Ack = true

					case ACK:
							// ACK has been recieved
							// If yours, do nothing. Ack msg will contain ack-senders ID
							// Else, add true to RecievedAck in matrix
							// If all RecievedAck is true, send to ordermanager
							// Else, wait for the rest or ??resend??
							// Set ack to false

						addSingleAck(AckMatrix, ResendMatrixAck, message)
						fmt.Println("Added true to someone else")

						if message.ID != elevator.ID {
							// Add true to RecievedAck with correct ID and order
						}

						message.Ack = false

							//if (message.ID != elevator.ID) {
								//fmt.Println("ID: ", elevator.ID, " recieved ack from: ", message.ID)
								//recievedAck := []Message{{Select: ACK, Done: false, ID: message.ID, Floor: message.Floor, Button: message.Button, State: message.State,
									//Dir: message.Dir, Ack: message.Ack, ResendMatrix: message.ResendMatrix, Matrix: message.Matrix}}
							//}

					case SendMatrix:
						MatrixUpdatech <- message
						message.Ack = true


					case UpdatedMatrix:
						MatrixUpdatech <- message
						message.Ack = true

					}
					
					if message.Ack && (message.SenderID != elevator.ID) {
						// Make an ACK msg, and send
						sendAck := []Message{{Select: ACK, Done: false, SenderID: elevator.ID, ID: message.ID, Floor: message.Floor, Button: message.Button, 
							State: message.State, Dir: message.Dir, Ack: message.Ack, ResendMatrix: message.ResendMatrix, Matrix: message.Matrix}}
						syncChans.ChangeInOrderch <- sendAck

						// If you have recieved an order, do it and set ack to true. Send ack back
					}
				}
			}*/

		
		case <- ackTicker.C:
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
					fmt.Println( newID, "is online")
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
					message := Message{Select: UpdateOffline, ID: newID}
					UpdateElevStatusch <- message
				} else {

					Online = false
					fmt.Println("I am offline")
				}
			}
		}
	}
}




func addAck (matrix [4+NumFloors][3*NumElevators]AckStruct, resendMatrixAck AckStruct, messages []Message, data int) {
	for _, message := range messages {

		for elev := 0; elev < NumElevators; elev++ {
			switch message.Select {

			case NewOrder:
				matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].Data = data
				matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].AwaitingAck[elev] = true
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
	}
}

func addSingleAck(matrix [4+NumFloors][3*NumElevators]AckStruct, resendMatrixAck AckStruct, message Message) {
	switch message.Select {
		case NewOrder:
			matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].AwaitingAck[message.SenderID] = false
			fmt.Println("Setting (", len(matrix) - message.Floor - 1, ", ", message.SenderID*3 + int(message.Button), ") to false")
			fmt.Println(message.ID*3)

		case OrderComplete:
			matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].AwaitingAck[message.SenderID] = false

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
func checkAck(matrix [4+NumFloors][3*NumElevators]AckStruct, resendMatrixAck AckStruct, message Message) bool{
	awaiting := false
	for elev := 0; elev < NumElevators; elev++ {
		switch message.Select {
		case NewOrder:
			if (matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].AwaitingAck[elev] == true) {
				fmt.Println("awaiting")
				awaiting = true
			}

		case OrderComplete:
			if (matrix[len(matrix) - message.Floor - 1][message.SenderID*3 + int(message.Button)].AwaitingAck[elev] == true) {
				awaiting = true
			}

		case UpdateStates:
			if (matrix[1][message.SenderID*3].AwaitingAck[elev] == true) {
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

