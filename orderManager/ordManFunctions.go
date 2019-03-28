package orderManager

import (
	. "../Config"
	"../IO"
)

func UpdateElevStatus(elevatorMatrix [][]int, LocalStateUpdatech chan Message,
	SyncUpdatech chan []Message, localElev Elevator) {
	for {
		select {
		case message := <- LocalStateUpdatech:
			InsertID(message.ID, elevatorMatrix)
			InsertState(message.ID, message.State, elevatorMatrix)
			InsertDirection(message.ID, message.Dir, elevatorMatrix)
			InsertFloor(message.ID, message.Floor, elevatorMatrix)

			OutMessage := []Message{message}

			OutMessage[0].SenderID = localElev.ID

			if !message.Done {
				SyncUpdatech <- OutMessage

			}
		}
	}
}

func checkLostOrders(elevatorMatrix [][]int, localElev Elevator, NewLocalOrderChan chan int) {
	for floor := 0; floor < NumFloors; floor++ {
		for button := 0; button < 3; button++ {
			for elev := 0; elev < NumElevators; elev++ {
				if elevatorMatrix[1][3*elev] == int(UNDEFINED) && elev != localElev.ID && button != 2 {
					// check others matrix
					if elevatorMatrix[len(elevatorMatrix)-floor-1][button+elev*NumElevators] == 1 {
						lostOrder := ButtonEvent{Floor: floor, Button: ButtonType(button)}
						addOrder(localElev.ID, elevatorMatrix, lostOrder)
						lightOrder := Message{ID: elev, Floor: floor, Button: ButtonType(button)}
						setLight(lightOrder, localElev)
						clearLostOrders(floor, elevatorMatrix, localElev, elev)
						NewLocalOrderChan <- floor
					}
				} else if elevatorMatrix[1][3*elev] != int(UNDEFINED) && elev == localElev.ID {
					if elevatorMatrix[len(elevatorMatrix)-floor-1][button+elev*NumElevators] == 1 {
						lostOrder := ButtonEvent{Floor: floor, Button: ButtonType(button)}
						lightOrder := Message{ID: elev, Floor: floor, Button: ButtonType(button)}
						setLight(lightOrder, localElev)
						addOrder(localElev.ID, elevatorMatrix, lostOrder)
						NewLocalOrderChan <- floor
					}
				}
			}
		}
	}
}

func updateOrdersInMatrix(newMatrix [][]int, oldMatrix [][]int, id int) [][]int {
	for i := 0; i < (4 + NumFloors); i++ {
		for j := 0; j < 3*NumElevators; j++ {
			if j != 3*id || i > 3 {
				newMatrix[i][j] = oldMatrix[i][j]
			}
		}
	}
	return newMatrix
}

func addOrder(id int, matrix [][]int, buttonPressed ButtonEvent) [][]int {
	matrix[NumFloors+3-buttonPressed.Floor][id*NumElevators+int(buttonPressed.Button)] = 1
	return matrix
}

func clearFloors(currentFloor int, elevatorMatrix [][]int, localElev Elevator, messageID int) {
	for button := 0; button < 3; button++ {
		elevatorMatrix[len(elevatorMatrix)-currentFloor-1][button+messageID*NumElevators] = 0
	}
}

func clearLostOrders(currentFloor int, elevatorMatrix [][]int, localElev Elevator, messageID int) {
	for button := 0; button < 2; button++ {
		elevatorMatrix[len(elevatorMatrix)-currentFloor-1][button+messageID*NumElevators] = 0
	}
	if messageID == localElev.ID {
		elevatorMatrix[len(elevatorMatrix)-currentFloor-1][2+messageID*NumElevators] = 0
	}
}

func InsertID(id int, matrix [][]int) {
	matrix[0][id*3] = id
}

func InsertFloor(id int, newFloor int, matrix [][]int) {
	matrix[2][id*3] = newFloor
}

func InsertState(id int, state int, matrix [][]int) {
	matrix[1][3*id] = state
}

func InsertDirection(id int, dir MotorDirection, matrix [][]int) {
	matrix[3][3*id] = int(dir)
}

func setLight(illuminateOrder Message, localElev Elevator) {
	if illuminateOrder.ID == localElev.ID {
		io.SetButtonLamp(illuminateOrder.Button, illuminateOrder.Floor, true)
	} else if illuminateOrder.Button == BT_Cab {

	} else {
		io.SetButtonLamp(illuminateOrder.Button, illuminateOrder.Floor, true)
	}
}

func clearLight(LocalOrderFinished int, localElev Elevator, messageID int) {
	if localElev.ID == messageID {
		io.SetButtonLamp(BT_Cab, LocalOrderFinished, false)
	}
	io.SetButtonLamp(BT_HallUp, LocalOrderFinished, false)
	io.SetButtonLamp(BT_HallDown, LocalOrderFinished, false)
}
