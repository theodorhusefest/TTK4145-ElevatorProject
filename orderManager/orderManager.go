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
  UpdateElevatorChan chan Message
  LocalOrderFinishedChan chan int
}


func OrderManager(OrderManagerChans OrderManagerChannels, NewGlobalOrderChan chan ButtonEvent, NewLocalOrderChan chan int,  elevatorMatrix [][]int, OutGoingMsg chan Message, ChangeInOrderch chan Message, elevatorConfig config.ElevConfig) {
  message := Message{
  }
  localOrder := ButtonEvent{
  }
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

      // Update message to be sent to everyone. Select = 1 for new order
      message.Select = 1
      message.Done = false
      message.ID = elevatorConfig.ElevID // ELEV-ID TO DEDICATED ELEVATOR
      message.Floor = NewGlobalOrder.Floor
      message.Button = NewGlobalOrder.Button

      // Send message to sync
      ChangeInOrderch <- message


      addOrder(elevatorConfig.ElevID, elevatorMatrix, NewGlobalOrder) // ELEV-ID TO DEDICATED ELEVATOR
      setLight(message, elevatorConfig)

      // if costfunction gives local elevator the order:
      NewLocalOrderChan <- int(NewGlobalOrder.Floor)

      // Print updated matrix for fun
      utilities.PrintMatrix(elevatorMatrix,4,3)


      default:

      //var newHallOrders []Message
      /*newHallOrders =*/ hallOrderAssigner.AssignHallOrder(NewGlobalOrder, elevatorMatrix)
      //fmt.Println("New Hall orders", newHallOrders)
      fmt.Println()
      // ???????????????????       Send new_order to everyonerderManager.InsertState(elevatorConfig.ElevID, 0, elevatorMa     ????

      // Update message to be sent to everyone. Select = 1 for new order
      message.Select = 1
      message.Done = false
      message.ID = elevatorConfig.ElevID // ELEV-ID TO DEDICATED ELEVATOR
      message.Floor = NewGlobalOrder.Floor
      message.Button = NewGlobalOrder.Button

      // Send message to sync
      ChangeInOrderch <- message

      // Wait for sync to say everyone knows the same

      //  Update local matrix, addOrder
      addOrder(elevatorConfig.ElevID, elevatorMatrix, NewGlobalOrder) // ELEV-ID TO DEDICATED ELEVATOR
      setLight(message, elevatorConfig)

      // if costfunction gives local elevator the order:
      NewLocalOrderChan <- int(NewGlobalOrder.Floor)

      // Print updated matrix for fun
      utilities.PrintMatrix(elevatorMatrix,4,3)

    }

    // -----------------------------------------------------------------------------------------------------Case triggered by elevator done with order
    case LocalOrderFinished := <- OrderManagerChans.LocalOrderFinishedChan:

      // Update message to be sent to everyone. Select = 2 for order done
      message.Select = 2
      message.Done = false
      message.ID = elevatorConfig.ElevID
      message.Floor = LocalOrderFinished

      // Send message to sync
      ChangeInOrderch <- message

      // Wait for sync to say everyone knows the same

      // Clear local matrix and lights
      clearFloors(LocalOrderFinished, elevatorMatrix, elevatorConfig.ElevID)
      clearLight(LocalOrderFinished)

      // Print updated matrix for fun
      utilities.PrintMatrix(elevatorMatrix,4,3)



    // -------------------------------------------------------------------------------------------------------Case triggered by incomming update (New_order, order_done etc.)
  case newUpdateFromSync := <- OrderManagerChans.UpdateElevatorChan:
      message = newUpdateFromSync
      if !message.Done {
        //SELECT = 1: NEW ORDER
        if message.Select == 1 {
          //NEW ORDER
          localOrder.Floor = message.Floor
          localOrder.Button = message.Button
          addOrder(message.ID, elevatorMatrix, localOrder)
          setLight(message, elevatorConfig)


          if message.ID == elevatorConfig.ElevID {
            //THIS ELEVATOR GOT THE JOB, TRIGGER FSM
          }
        }
        // SELECT = 2: AN ORDER IS FINISHED
        if message.Select == 2 {
          clearFloors(message.Floor, elevatorMatrix, message.ID)
          clearLight(message.Floor)
        }
        // SELECT = 3: NEW CHANGE IN STATE/FLOOR/DIR FOR AN ELEVATOR
        if message.Select == 3 {
          InsertState(message.ID, message.State, elevatorMatrix)
          // FIKS!!!!!!!!!  InsertDirection(message.ID, message.Dir, elevatorMatrix)
          // LAG FUNKSJON SOM SETTER INN FLOOR ????????????
        }
      message.Done = true
      }




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

func InsertState(id int, state int, matrix [][]int){
  matrix[1][3*id] = state
}

func InsertDirection(id int, elevator config.Elevator, matrix [][]int){
  switch elevator.Dir{
    case DIR_Up:
      matrix[3][3*id] = 1
    case DIR_Down:
      matrix[3][3*id] = -1
    case DIR_Stop:
      matrix[3][3*id] = 0
  }
}


/*
func setLight(newGlobalOrder ButtonEvent) {
  //Set button lamp
  io.SetButtonLamp(newGlobalOrder.Button, newGlobalOrder.Floor, true)
}
*/

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
