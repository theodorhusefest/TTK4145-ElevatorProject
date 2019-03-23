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

	AckMatrix[0][0].Data = 1

	broadCastTicker := time.NewTicker(100 * time.Millisecond)
	ackTicker := time.NewTicker(100 * time.Millisecond)
	for {
		select {

		// --------------------------------------------------------------------------Case triggered by local ordermanager, change in order
		case changeInOrder := <-syncChans.ChangeInOrderch:
			//Håndter endring som kom fra ordermanager. Send alt inn på message og sett message.Done = false
			//message := changeInOrder
      //utilities.PrintMatrix(elevatorMatrix, 4,3)

      if Online {
				addOrderToAck(AckMatrix, ResendMatrixAck, changeInOrder)
				select {
				case <-broadCastTicker.C:
					//fmt.Println(elevator.ID, "is sending outgoing message ")
					syncChans.OutGoingMsg <- changeInOrder
				}
			

			} else {
				for _, message := range changeInOrder {
					if !(message.Done) {
						//SELECT = 1: NEW ORDER
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
							// ACKNOWLEDGE

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

					//fmt.Println(elevator.ID, "is recieving incomming message Type:", MessageType(message.Select), "from", message.ID)
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
						fmt.Println("Ack recieved by: ", message.ID)

					case SendMatrix:
						MatrixUpdatech <- message
						message.Ack = true


					case UpdatedMatrix:
						MatrixUpdatech <- message
						message.Ack = true
	
					}
					if message.Ack {
						sendAck := []Message{{Select: ACK, Done: message.Done, ID: message.ID, Floor: message.Floor, Button: message.Button, State: message.State, 
							Dir: message.Dir, Ack: message.Ack, ResendMatrix: message.ResendMatrix, Matrix: message.Matrix}}
						fmt.Println("Sending Ack")
						//syncChans.OutGoingMsg <- sendAck // Lag til en array
						syncChans.ChangeInOrderch <- sendAck
					}
				}
			}

	
		case <- ackTicker.C:





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




func addOrderToAck (matrix [4+NumFloors][3*NumElevators]AckStruct, resendMatrixAck AckStruct, messages []Message) {
	for _, message := range messages {

		for elev := 0; elev < NumElevators; elev++ {
			switch message.Select {

			case NewOrder:
				matrix[len(matrix) - message.Floor - 1][message.ID*3 + int(message.Button)].Data = 1
				matrix[len(matrix) - message.Floor - 1][message.ID*3 + int(message.Button)].AwaitingAck[elev] = true
				matrix[len(matrix) - message.Floor - 1][message.ID*3 + int(message.Button)].RecievedAck[elev] = false

			case OrderComplete:
				for i := 0; i < 3; i++ {
					matrix[len(matrix) - message.Floor - 1][message.ID*3 + i].Data = 0
					matrix[len(matrix) - message.Floor - 1][message.ID*3 + i].AwaitingAck[elev] = true
					matrix[len(matrix) - message.Floor - 1][message.ID*3 + i].RecievedAck[elev] = false
				} 

			case UpdateStates:
				matrix[1][message.ID*3].Data = message.State
				matrix[1][message.ID*3].AwaitingAck[elev] = true
				matrix[1][message.ID*3].RecievedAck[elev] = false

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

func addAck(matrix [4+NumFloors][3*NumElevators]AckStruct, resendMatrixAck AckStruct, messages []Message) {

}