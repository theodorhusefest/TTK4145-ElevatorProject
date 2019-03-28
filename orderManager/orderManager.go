package orderManager

import (
	. "../Config"
	"../IO"
	"../hallRequestAssigner"
	"fmt"
	"time"
)

/*
Matrix  ID_1        -----------   -----     ID_2        -----------   -----     ID_3        -----------   -----
        State_1     -----------   -----     State_2     -----------   -----     State_3     -----------   -----
        FLoor_1     -----------   -----     FLoor_2     -----------   -----     FLoor_3     -----------   -----
        Dir_1       -----------   -----     Dir_2       -----------   -----     Dir_3       -----------   -----
        --------    Down_floor_4  Cab_4     --------    Down_floor_4  Cab_4     --------    Down_floor_4  Cab_4
        Up_floor_3  Down_floor 3  Cab_3     Up_floor_3  Down_floor 3  Cab_3     Up_floor_3  Down_floor 3  Cab_3
        Up_floor_2  Down_floor_2  Cab_2     Up_floor_2  Down_floor_2  Cab_2     Up_floor_2  Down_floor_2  Cab_2
        Up_floor_1  --------      Cab_1     Up_floor_1  --------      Cab_1     Up_floor_1  --------      Cab_1
*/

type OrderManagerChannels struct {
	LocalOrderFinishedChan chan int
	NewLocalOrderch        chan Message
	UpdateOrderch          chan Message
	MatrixUpdatech         chan Message
}

func OrderManager(	elevatorMatrix [][]int, localElev Elevator, 
					OrderManagerChans OrderManagerChannels,
					ButtonPressedchn chan ButtonEvent, 
					NewLocalOrderChan chan int, 
					SyncUpdatech chan []Message, 
					UpdateElevStatusch chan Message, 
					GlobalStateUpdatech chan Message) {

	OrderTimedOut := time.NewTimer(5 * time.Second)

	for {
		select {


		// -----------------------------------------------------------------------------------------------------Case triggered by local button
		case ButtonPressed := <-ButtonPressedchn:

			if elevatorMatrix[1][localElev.ID*NumElevators] != 3 {
				switch ButtonPressed.Button {

				case BT_Cab:

					newCabOrder := []Message{{Select: NewOrder, Done: false, SenderID: localElev.ID, 
												ID: localElev.ID, Floor: ButtonPressed.Floor, 
												Button: ButtonPressed.Button}}
					// Send message to sync
					fmt.Println("Cab order at elevator:", localElev.ID)
					SyncUpdatech <- newCabOrder

				default:

					newHallOrders := hallOrderAssigner.AssignHallOrder(ButtonPressed, elevatorMatrix, localElev)


					fmt.Println("Hall order at elevator:", localElev.ID)

					SyncUpdatech <- newHallOrders
				}

			OrderTimedOut.Reset(10 * time.Second)

			}
			

		case OrderUpdate := <-OrderManagerChans.UpdateOrderch:
			switch OrderUpdate.Select {
			case NewOrder:
				localOrder := ButtonEvent{Floor: OrderUpdate.Floor, Button: OrderUpdate.Button}
				addOrder(OrderUpdate.ID, elevatorMatrix, localOrder)
				setLight(OrderUpdate, localElev)
				if OrderUpdate.ID == localElev.ID {
					NewLocalOrderChan <- OrderUpdate.Floor
				}


			case OrderComplete:
				clearFloors(OrderUpdate.Floor, elevatorMatrix, localElev, OrderUpdate.ID)
				clearLight(OrderUpdate.Floor, localElev, OrderUpdate.ID)
			}
			OrderTimedOut.Reset(10 * time.Second)

		case StateUpdate := <-GlobalStateUpdatech:

			switch StateUpdate.Select {
			case UpdateStates:

				InsertID(StateUpdate.ID, elevatorMatrix)
				InsertState(StateUpdate.ID, StateUpdate.State, elevatorMatrix)
				InsertDirection(StateUpdate.ID, StateUpdate.Dir, elevatorMatrix)
				InsertFloor(StateUpdate.ID, StateUpdate.Floor, elevatorMatrix)

			case UpdateOffline:
				InsertState(StateUpdate.ID, int(UNDEFINED), elevatorMatrix)

			}

		case MatrixUpdate := <-OrderManagerChans.MatrixUpdatech:

			switch MatrixUpdate.Select {
			case SendMatrix:
				fmt.Println("Sending matrix")
				outMessage := []Message{{Select: UpdatedMatrix, SenderID: localElev.ID, Matrix: elevatorMatrix, ID: MatrixUpdate.ID}}
				SyncUpdatech <- outMessage

			case UpdatedMatrix:
				if MatrixUpdate.ID == localElev.ID {
					elevatorMatrix = updateOrdersInMatrix(elevatorMatrix, MatrixUpdate.Matrix, MatrixUpdate.ID)
					outMessage := []Message{{Select:UpdateStates ,ID: localElev.ID, State: int(localElev.State), Floor: elevatorMatrix[2][localElev.ID*3], Dir: localElev.Dir}}
					SyncUpdatech <- outMessage
				}
			}

		// -----------------------------------------------------------------------------------------------------Case triggered by elevator done with order
		case LocalOrderFinished := <-OrderManagerChans.LocalOrderFinishedChan:

			outMessage := []Message{{Select: OrderComplete, SenderID: localElev.ID, Done: false, ID: localElev.ID, Floor: LocalOrderFinished}}

			SyncUpdatech <- outMessage


		// ------------------------------------------------------------------------------------------------------- Case triggered every 5 seconds to check if orders left
		case <- OrderTimedOut.C:
			fmt.Println("Checking for lost orders")
      		checkLostOrders(elevatorMatrix, localElev, NewLocalOrderChan)
			OrderTimedOut.Reset(10 * time.Second)



			// -------------------------------------------------------------------------------------------------------Case triggered by incomming update (New_order, order_done etc.)
		}
	}
}

func UpdateElevStatus(elevatorMatrix [][]int, UpdateElevStatusch chan Message, SyncUpdatech chan []Message, localElev Elevator) {
  for {
    select {
    case message := <-UpdateElevStatusch:
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
            if elevatorMatrix[len(elevatorMatrix)-floor - 1][button+elev*NumElevators] == 1 {
              lostOrder := ButtonEvent{Floor: floor, Button: ButtonType(button)}
              addOrder(localElev.ID, elevatorMatrix, lostOrder)
              lightOrder := Message{ID:elev, Floor: floor, Button: ButtonType(button)}
              setLight(lightOrder, localElev)
              clearFloors2(floor, elevatorMatrix, localElev, elev)
              NewLocalOrderChan <- floor
            }
        } else if elevatorMatrix[1][3*elev] != int(UNDEFINED) && elev == localElev.ID  {
            if elevatorMatrix[len(elevatorMatrix)-floor - 1][button+elev*NumElevators] == 1 {
              lostOrder := ButtonEvent{Floor: floor, Button: ButtonType(button)}
              lightOrder := Message{ID:elev, Floor: floor, Button: ButtonType(button)}
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

func clearFloors2(currentFloor int, elevatorMatrix [][]int, localElev Elevator, messageID int) {
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
	switch dir {
	case DIR_Up:
		matrix[3][3*id] = 1
	case DIR_Down:
		matrix[3][3*id] = -1
	case DIR_Stop:
		matrix[3][3*id] = 0
	}
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
