package orderManager

import(
  . "../Config"
  "../Utilities"
  "../IO"
  "../Config"
  //"fmt"
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
  UpdateElevatorChan chan Elevator
  LocalOrderFinishedChan chan int
}


func OrderManager(OrderManagerChans OrderManagerChannels, NewGlobalOrderChan chan ButtonEvent, NewLocalOrderChan chan int,  elevatorMatrix [][]int, OutGoingOrder chan Message, MessageToSend chan Message, MessageRecieved chan Message, elevatorConfig config.ElevConfig) {
  message := Message{
  }
  localOrder := ButtonEvent{
  }
  for {
    select {
    case newGlobalOrder := <- NewGlobalOrderChan:

      // Costfunction(elevatorMatrix)


      // Update matrix
      addOrder(elevatorConfig.ElevID, elevatorMatrix, newGlobalOrder)

      // Send to network
      message.ID = 1
      message.Floor = newGlobalOrder.Floor
      message.Button = newGlobalOrder.Button
      MessageToSend <- message


      // if own elevator send to newLocalOrder
      NewLocalOrderChan <- int(newGlobalOrder.Floor)
      setLight(newGlobalOrder)
    //case UpdateElevator := <- UpdateElevatorChan:

    case LocalOrderFinished := <- OrderManagerChans.LocalOrderFinishedChan:
      clearFloors(LocalOrderFinished, elevatorMatrix, elevatorConfig.ElevID)
      clearLight(LocalOrderFinished)
      utilities.PrintMatrix(elevatorMatrix,4,3)

    case newNetworkOrder := <- MessageRecieved:
      localOrder.Floor = newNetworkOrder.Floor
      localOrder.Button = newNetworkOrder.Button
      addOrder(newNetworkOrder.ID, elevatorMatrix, localOrder)




    // case orderFinished

    }
  }
}




func addOrder(id int, matrix [][]int, buttonPressed ButtonEvent) [][]int{
  matrix[NumFloors+3-buttonPressed.Floor][id*NumElevators + int(buttonPressed.Button)] = 1
  return matrix
}


func clearFloors(currentFloor int, elevatorMatrix [][]int, id int) {
	for button:=0; button < NumFloors; button++ {
		elevatorMatrix[len(elevatorMatrix)-currentFloor-1][button+id*NumElevators] = 0
	}
}

func setLight(newGlobalOrder ButtonEvent) {
  //Set button lamp
  io.SetButtonLamp(newGlobalOrder.Button, newGlobalOrder.Floor, true)
}

func clearLight(LocalOrderFinished int) {
	io.SetButtonLamp(BT_Cab, LocalOrderFinished, false)
  io.SetButtonLamp(BT_HallUp, LocalOrderFinished, false)
  io.SetButtonLamp(BT_HallDown, LocalOrderFinished, false)
}
