package orderManager

import(
  . "../Config"
  "../Utilities"
  "../IO"
  "../Config"
  "../hallRequestAssigner"
  "fmt"
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



type OrderManagerChannels struct{
  UpdateElevatorChan chan []Message
  LocalOrderFinishedChan chan int
  UpdateElevStatusch chan []Message
  PrintLiveMatrixCh chan string
}


func OrderManager(OrderManagerChans OrderManagerChannels, NewGlobalOrderChan chan ButtonEvent, NewLocalOrderChan chan int,
        elevatorMatrix [][]int, OutGoingMsg chan []Message, ChangeInOrderch chan []Message, SendFullMatrixch chan [][]int ,elevatorConfig config.ElevConfig) {
  localOrder := ButtonEvent{}
  for {
    select {

      /*
      1: Ordre tas imot av en heis.
      2: Den heisen kjører kostfunksjon og bestemmer hvem som får jobben.
      3: Heisen sender ordren til alle andre heiser, så alle er oppdatert.
      4: Alle skrur
      5: Den heisen som får jobben, trigger sin egen FSM med        NewLocalOrderChan <- int(newGlobalOrder.Floor)
      6: Heisen som har utført et oppdrag oppdaterer de andre, så alle vet det samme
      7: Alle fjerner ordren lokalt og evt. skrur av lys
      */

    // -----------------------------------------------------------------------------------------------------Case triggered by local button
    case NewGlobalOrder := <- NewGlobalOrderChan:

      switch NewGlobalOrder.Button {

      case BT_Cab:

      outMessage := []Message {{Select: 1, Done: false, ID: elevatorConfig.ElevID, Floor: NewGlobalOrder.Floor, Button: NewGlobalOrder.Button}}
      // Send message to sync

      ChangeInOrderch <- outMessage

     // addOrder(elevatorConfig.ElevID, elevatorMatrix, NewGlobalOrder) // ELEV-ID TO DEDICATED ELEVATOR
    //  setLight(outMessage[0], elevatorConfig)


      default:

      //var newHallOrders []Message
      newHallOrders := hallOrderAssigner.AssignHallOrder(NewGlobalOrder, elevatorMatrix)
      //fmt.Println("New Hall orders", newHallOrders)
      fmt.Println()
      // ???????????????????       Send new_order to everyonerderManager.InsertState(elevatorConfig.ElevID, 0, elevatorMa     ????

      // Send message to sync                //time.Sleep(10*time.Second)

      ChangeInOrderch <- newHallOrders
      /*
      // Wait for sync to say everyone knows the same

      //  Update local matrix, addOrder
      addOrder(elevatorConfig.ElevID, elevatorMatrix, NewGlobalOrder) // ELEV-ID TO DEDICATED ELEVATOR
      setLight(message, elevatorConfig)

      // if costfunction gives local elevator the order:
      NewLocalOrderChan <- int(NewGlobalOrder.Floor)

      // Print updated matrix for fun
      utilities.PrintMatrix(elevatorMatrix,4,3)
*/
    }

    // -----------------------------------------------------------------------------------------------------Case triggered by elevator done with order
    case LocalOrderFinished := <- OrderManagerChans.LocalOrderFinishedChan:

      // Update message to be sent to everyone. Select = 2 for order done

      outMessage := []Message {{Select: 2, Done: false, ID: elevatorConfig.ElevID, Floor: LocalOrderFinished}}

      // Send message to sync
      ChangeInOrderch <- outMessage

      // Wait for sync to say everyone knows the same

      // Print updated matrix for fun)



//  case elevStatus := <-OrderManagerChans.UpdateElevStatus:
//      ChangeInOrderch <- elevStatus
      //fmt.Println(elevStatus)


    // -------------------------------------------------------------------------------------------------------Case triggered by incomming update (New_order, order_done etc.)
  case newUpdateFromSync := <- OrderManagerChans.UpdateElevatorChan:
      inMessages := newUpdateFromSync
      for _, message := range inMessages {
          if !(message.Done) {
            //SELECT = 1: NEW ORDER

            switch  message.Select {
            case 1:
                localOrder.Floor = message.Floor
                localOrder.Button = message.Button
                fmt.Println("MessageID", message.ID, localOrder)
                addOrder(message.ID, elevatorMatrix, localOrder)
                setLight(message, elevatorConfig)
                utilities.PrintMatrix(elevatorMatrix, elevatorConfig.NumFloors,elevatorConfig.NumElevators)



                if message.ID == elevatorConfig.ElevID {
                    NewLocalOrderChan <- message.Floor
                    fmt.Println("Sending LocalOrder")
                }

            case 2:
                clearFloors(message.Floor, elevatorMatrix, message.ID)
                clearLight(message.Floor)

            case 3:
                InsertID(message.ID, elevatorMatrix)
                InsertState(message.ID, message.State, elevatorMatrix)
                InsertDirection(message.ID, message.Dir, elevatorMatrix)
                InsertFloor(message.ID, message.Floor, elevatorMatrix)


            case 4:
                // ACKNOWLEDGE
            case 5:
                fmt.Println("Someone new on network, sending matrix")
                outMessage := []Message{{Select: 6, Matrix: elevatorMatrix, ID: message.ID}}
                ChangeInOrderch <- outMessage


            case 6:
                if message.ID == elevatorConfig.ElevID{
                    fmt.Println("Resetting matrix")
                    elevatorMatrix = updateOrdersInMatrix(elevatorMatrix, message.Matrix, message.ID)
                }
                utilities.PrintMatrix(elevatorMatrix, elevatorConfig.NumFloors,elevatorConfig.NumElevators)

            case 7:
                InsertState(message.ID, int(OFFLINE), elevatorMatrix)

            }
            message.Done = true
        }
      }
      fmt.Println("End for for loop")
    }
  }
}


func PrintLiveMatrix(input chan string) {
  for {
    select {
    case <- input:
      utilities.PrintMatrix(elevatorMatrix, elevatorConfig.NumFloors,elevatorConfig.NumElevators)

    }
  }
}



func updateOrdersInMatrix(newMatrix [][]int, oldMatrix [][]int, id int) [][]int {
    for i := 0; i < (4+NumFloors); i ++ {
        for j := 0; j < 3*NumElevators; j ++ {
            if j != 3*id || i > 3 {
                newMatrix[i][j] = oldMatrix[i][j]
            }

        }
    }
    return newMatrix
}


func UpdateElevStatus(UpdateElevStatusch chan []Message, ChangeInOrderch chan []Message, elevatorMatrix [][]int) {
    for {
        select {
        case inMessages := <- UpdateElevStatusch:
            for _, message := range inMessages {
                InsertID(message.ID, elevatorMatrix)
                InsertState(message.ID, message.State, elevatorMatrix)
                InsertDirection(message.ID, message.Dir, elevatorMatrix)
                InsertFloor(message.ID, message.Floor, elevatorMatrix)

            }
            ChangeInOrderch <- inMessages
        }

    }
}





func addOrder(id int, matrix [][]int, buttonPressed ButtonEvent) [][]int{
  matrix[NumFloors+3-buttonPressed.Floor][id*NumElevators + int(buttonPressed.Button)] = 1
  return matrix
}


func clearFloors(currentFloor int, elevatorMatrix [][]int, id int) {
	for button:=0; button < NumElevators; button++ {
		elevatorMatrix[len(elevatorMatrix)-currentFloor-1][button+id*NumElevators] = 0
	}
}

func InsertID(id int, matrix [][]int) {
    matrix[0][id*3] = id
}


func InsertFloor(id int, newFloor int, matrix [][]int){
  matrix[2][id*3] = newFloor
}

func InsertState(id int, state int, matrix [][]int){
  matrix[1][3*id] = state
}

func InsertDirection(id int, dir MotorDirection, matrix [][]int){
  switch dir{
    case DIR_Up:
      matrix[3][3*id] = 1
    case DIR_Down:
      matrix[3][3*id] = -1
    case DIR_Stop:
      matrix[3][3*id] = 0
  }
}



func setLight(illuminateOrder Message, elevatorConfig config.ElevConfig) {
  if illuminateOrder.ID == elevatorConfig.ElevID {
    io.SetButtonLamp(illuminateOrder.Button, illuminateOrder.Floor, true)
  } else if int(illuminateOrder.Button) == 2 {

  } else {
    io.SetButtonLamp(illuminateOrder.Button, illuminateOrder.Floor, true)
  }
}


func clearLight(LocalOrderFinished int) {
	io.SetButtonLamp(BT_Cab, LocalOrderFinished, false)
  io.SetButtonLamp(BT_HallUp, LocalOrderFinished, false)
  io.SetButtonLamp(BT_HallDown, LocalOrderFinished, false)
}




// FIKS
// Trigg fsm hvis ordre ikke blir gjort
// Button cab light
// initialize lights
// Ack
// Bug med 2 knapper samtidig
//
