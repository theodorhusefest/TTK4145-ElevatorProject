package orderManager

import (
	. "../Config"
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

func OrderManager(elevatorMatrix [][]int, localElev Elevator,
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
				outMessage := []Message{{Select: UpdatedMatrix, SenderID: localElev.ID,
					Matrix: elevatorMatrix, ID: MatrixUpdate.ID}}
				SyncUpdatech <- outMessage

			case UpdatedMatrix:
				if MatrixUpdate.ID == localElev.ID {
					elevatorMatrix = updateOrdersInMatrix(elevatorMatrix, MatrixUpdate.Matrix, MatrixUpdate.ID)
					outMessage := []Message{{Select: UpdateStates, ID: localElev.ID,
						State: int(localElev.State), Floor: elevatorMatrix[2][localElev.ID*3],
						Dir: localElev.Dir}}
					SyncUpdatech <- outMessage
				}
			}

		case LocalOrderFinished := <-OrderManagerChans.LocalOrderFinishedChan:

			outMessage := []Message{{Select: OrderComplete, SenderID: localElev.ID, Done: false, ID: localElev.ID, Floor: LocalOrderFinished}}

			SyncUpdatech <- outMessage

		case <-OrderTimedOut.C:
			fmt.Println("Checking for lost orders")
			checkLostOrders(elevatorMatrix, localElev, NewLocalOrderChan)
			OrderTimedOut.Reset(10 * time.Second)

		}
	}
}
