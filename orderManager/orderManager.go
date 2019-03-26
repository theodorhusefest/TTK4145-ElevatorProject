package orderManager

import (
	. "../Config"
	//"../Utilities"
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

func OrderManager(	elevatorMatrix [][]int, elevator Elevator, OrderManagerChans OrderManagerChannels,
					ButtonPressedchn chan ButtonEvent, NewLocalOrderChan chan int, ChangeInOrderch chan []Message, 
					UpdateElevStatusch chan Message, GlobalStateUpdatech chan Message) {

	OrderTimedOut := time.NewTimer(10 * time.Second)

	for {
		select {


		// -----------------------------------------------------------------------------------------------------Case triggered by local button
		case ButtonPressed := <-ButtonPressedchn:


			switch ButtonPressed.Button {

			case BT_Cab:

				newCabOrder := []Message{{Select: NewOrder, Done: false, SenderID: elevator.ID, ID: elevator.ID, Floor: ButtonPressed.Floor, Button: ButtonPressed.Button}}
				// Send message to sync
				fmt.Println("NewCabOrder = ", newCabOrder)
				ChangeInOrderch <- newCabOrder

			default:

				newHallOrders := hallOrderAssigner.AssignHallOrder(ButtonPressed, elevatorMatrix, elevator)


				fmt.Println("NewHallOrder = ", newHallOrders)

				ChangeInOrderch <- newHallOrders
			}
			OrderTimedOut.Reset(10 * time.Second)

		case OrderUpdate := <-OrderManagerChans.UpdateOrderch:
			switch OrderUpdate.Select {
			case NewOrder:
				localOrder := ButtonEvent{Floor: OrderUpdate.Floor, Button: OrderUpdate.Button}
				addOrder(OrderUpdate.ID, elevatorMatrix, localOrder)
				setLight(OrderUpdate, elevator)
				if OrderUpdate.ID == elevator.ID {
					NewLocalOrderChan <- OrderUpdate.Floor
				}


			case OrderComplete:
				clearFloors(OrderUpdate.Floor, elevatorMatrix, OrderUpdate.ID)
				clearLight(OrderUpdate.Floor, elevator, OrderUpdate.ID)
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
				fmt.Println("Resending matrix")
				outMessage := []Message{{Select: UpdatedMatrix, SenderID: elevator.ID, Matrix: elevatorMatrix, ID: MatrixUpdate.ID}}
				ChangeInOrderch <- outMessage

			case UpdatedMatrix:
				if MatrixUpdate.ID == elevator.ID {
					fmt.Println("Resetting matrix")
					elevatorMatrix = updateOrdersInMatrix(elevatorMatrix, MatrixUpdate.Matrix, MatrixUpdate.ID)
					outMessage := []Message{{Select:UpdateStates ,ID: elevator.ID, State: int(elevator.State), Floor: elevatorMatrix[2][elevator.ID*3], Dir: elevator.Dir}}
					ChangeInOrderch <- outMessage
				}
			}

		// -----------------------------------------------------------------------------------------------------Case triggered by elevator done with order
		case LocalOrderFinished := <-OrderManagerChans.LocalOrderFinishedChan:

			outMessage := []Message{{Select: OrderComplete, SenderID: elevator.ID, Done: false, ID: elevator.ID, Floor: LocalOrderFinished}}

			ChangeInOrderch <- outMessage


		// ------------------------------------------------------------------------------------------------------- Case triggered every 5 seconds to check if orders left
		case <- OrderTimedOut.C:
			fmt.Println("Checking for lost orders")
      		checkLostOrders(elevatorMatrix, elevator, NewLocalOrderChan)
			OrderTimedOut.Reset(10 * time.Second)



			// -------------------------------------------------------------------------------------------------------Case triggered by incomming update (New_order, order_done etc.)
		}
	}
}

func UpdateElevStatus(elevatorMatrix [][]int, UpdateElevStatusch chan Message, ChangeInOrderch chan []Message, elevator Elevator) {
  for {
    select {
    case message := <-UpdateElevStatusch:
      InsertID(message.ID, elevatorMatrix)
      InsertState(message.ID, message.State, elevatorMatrix)
      InsertDirection(message.ID, message.Dir, elevatorMatrix)
      InsertFloor(message.ID, message.Floor, elevatorMatrix)

      OutMessage := []Message{message}

      OutMessage[0].SenderID = elevator.ID

      if !message.Done {
        ChangeInOrderch <- OutMessage

      }
    }
  }
}

func checkLostOrders(elevatorMatrix [][]int, elevator Elevator, NewLocalOrderChan chan int) {
	for floor := 0; floor < NumFloors; floor++ {
		for button := 0; button < 3; button++ {
      for elev := 0; elev < NumElevators; elev++ {
        if elevatorMatrix[1][3*elev] == int(UNDEFINED) && elev != elevator.ID && button != 2 {
            // check others matrix
            if elevatorMatrix[len(elevatorMatrix)-floor - 1][button+elev*NumElevators] == 1 {
              lostOrder := ButtonEvent{Floor: floor, Button: ButtonType(button)}
              addOrder(elevator.ID, elevatorMatrix, lostOrder)
              clearFloors(floor, elevatorMatrix, elev)
              NewLocalOrderChan <- floor
            }
        } else if elevatorMatrix[1][3*elev] != int(UNDEFINED) && elev == elevator.ID  {
            if elevatorMatrix[len(elevatorMatrix)-floor - 1][button+elev*NumElevators] == 1 {
              lostOrder := ButtonEvent{Floor: floor, Button: ButtonType(button)}
              addOrder(elevator.ID, elevatorMatrix, lostOrder)
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

func clearFloors(currentFloor int, elevatorMatrix [][]int, id int) {
	for button := 0; button < NumElevators; button++ {
		elevatorMatrix[len(elevatorMatrix)-currentFloor-1][button+id*NumElevators] = 0
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

func setLight(illuminateOrder Message, elevator Elevator) {
	if illuminateOrder.ID == elevator.ID {
		io.SetButtonLamp(illuminateOrder.Button, illuminateOrder.Floor, true)
	} else if illuminateOrder.Button == BT_Cab {

	} else {
		io.SetButtonLamp(illuminateOrder.Button, illuminateOrder.Floor, true)
	}
}

func clearLight(LocalOrderFinished int, elevator Elevator, messageID int) {
	if elevator.ID == messageID {
		io.SetButtonLamp(BT_Cab, LocalOrderFinished, false)
	}
	io.SetButtonLamp(BT_HallUp, LocalOrderFinished, false)
	io.SetButtonLamp(BT_HallDown, LocalOrderFinished, false)
}
