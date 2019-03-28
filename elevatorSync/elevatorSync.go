package syncElevator

import (
	. "../Config"
	"../Network/network/peers"
	"fmt"
	"strconv"
	"time"
)

type SyncElevatorChannels struct {
	OutGoingMsg    chan []Message
	InCommingMsg   chan []Message
	SyncUpdatech   chan []Message
	PeerUpdate     chan peers.PeerUpdate
	TransmitEnable chan bool
}

func SyncElevator(elevatorMatrix [][]int, localElev Elevator,
	syncChans SyncElevatorChannels, UpdateOrderch chan Message,
	LocalStateUpdatech chan Message, GlobalStateUpdatech chan Message,
	MatrixUpdatech chan Message) {

	Online := false
	AckMatrix := [4 + NumFloors][3 * NumElevators]AckStruct{}
	ResendMatrixAck := AckStruct{}

	broadCastTicker := time.NewTicker(10 * time.Millisecond)
	ackTicker := time.NewTicker(500 * time.Millisecond)

	for {
		select {

		case syncUpdate := <-syncChans.SyncUpdatech:

			if Online {
				for _, message := range syncUpdate {
					// If the message is not an ACK, we expect to recieve from other, so set AwaitingAck to true
					if !message.Ack {
						for elev := 0; elev < NumElevators; elev++ {
							switch message.Select {
							case NewOrder:
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].Data = 1
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] = true
							case OrderComplete:
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].Data = 0
								AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] = true

							case UpdatedMatrix:
								ResendMatrixAck.AwaitingAck[elev] = true
							}
						}
					}
				}
				select {
				case <-broadCastTicker.C:
					syncChans.OutGoingMsg <- syncUpdate
				}

			} else { // If currently offline, send directly back to ordermanager
				for _, message := range syncUpdate {
					if !(message.Done) {
						message.Done = true
						switch message.Select {

						case NewOrder, OrderComplete:
							UpdateOrderch <- message

						case UpdateStates, UpdateOffline:
							LocalStateUpdatech <- message

						case ACK:

						case SendMatrix, UpdatedMatrix:
							MatrixUpdatech <- message

						}
					}
				}
			}

		case msgRecieved := <-syncChans.InCommingMsg:
			for _, message := range msgRecieved {
				if !(message.Done) {

					message.Done = true
					sendAck := Message{Select: message.Select, Done: false, SenderID: localElev.ID,
						ID: message.ID, Floor: message.Floor, Button: message.Button,
						State: message.State, Dir: message.Dir, Ack: false,
						ResendMatrix: message.ResendMatrix, Matrix: message.Matrix}

					// If the message is an ACK or from yourself, set AwaitingAck to false
					// If all ACKs have arrived, exceute order
					// If message is a new update from someone else, execute order and send ACK
					switch message.Select {
					case NewOrder, OrderComplete:
						if message.Ack || message.SenderID == localElev.ID {
							sendAck.Ack = false
							AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[message.SenderID] = false
							allReady := true
							for elev := 0; elev < NumElevators; elev++ {
								if AckMatrix[len(AckMatrix)-message.Floor-1][message.ID*3+int(message.Button)].AwaitingAck[elev] == true &&
									elevatorMatrix[1][elev*3] != 3 {
									allReady = false
								}
							}

							if allReady {
								UpdateOrderch <- message
							}

						} else if !message.Ack && message.SenderID != localElev.ID {
							sendAck.Ack = true
							UpdateOrderch <- message
						}

					case UpdateStates:
						sendAck.Ack = false
						GlobalStateUpdatech <- message

					case SendMatrix, UpdatedMatrix:
						if message.Ack || message.SenderID == localElev.ID {
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
						} else if !message.Ack || message.SenderID != localElev.ID {
							sendAck.Ack = true
							MatrixUpdatech <- message
						}
					}
					if sendAck.Ack {
						sendAck := []Message{sendAck}
						syncChans.SyncUpdatech <- sendAck
					}
				}
			}

		// At every tick, go through the AckMatix and resend message if still awaiting ACK
		case <- ackTicker.C:
			for floor := 4; floor < 4+NumFloors; floor++ {
				for button := 0; button < 3*NumElevators; button++ {
					for elev := 0; elev < NumElevators; elev++ {
						if AckMatrix[floor][button].AwaitingAck[elev] == true && elevatorMatrix[1][elev*3] != 3 && localElev.ID != elev {
							ResendOrder := []Message{}
							if AckMatrix[floor][button].Data == 1 {
								ResendOrder = []Message{{Select: NewOrder, Done: false,
									SenderID: localElev.ID, ID: (button / NumElevators),
									Floor: 7 - floor, Button: ButtonType(button % 3)}}
							} else {
								ResendOrder = []Message{{Select: OrderComplete, Done: false,
									SenderID: localElev.ID, ID: (button / NumElevators),
									Floor: 7 - floor, Button: ButtonType(button % 3)}}
							}
							syncChans.SyncUpdatech <- ResendOrder
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

		case p := <-syncChans.PeerUpdate:

			time.Sleep(500 * time.Millisecond)
			AckMatrix = [4 + NumFloors][3 * NumElevators]AckStruct{} // Reset AckMatrix when losing or getting a peer

			if len(p.New) > 0 {
				newID, _ := strconv.Atoi(p.New)
				if newID == localElev.ID && len(p.Peers) == 1 {
					Online = true
					fmt.Println(newID, "is online")
				} else if newID == localElev.ID && len(p.Peers) > 1 {

					Online = true
					fmt.Println(newID, "is online")

				} else if newID != localElev.ID && Online {
					message := Message{Select: SendMatrix, ID: newID}
					MatrixUpdatech <- message
				}
			}

			for _, peerLost := range p.Lost {
				newID, _ := strconv.Atoi(peerLost)
				if newID != localElev.ID {
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
