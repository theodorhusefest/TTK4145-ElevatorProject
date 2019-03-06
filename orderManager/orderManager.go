package orderManager

import(
  "../IO"
  //"../Initialize"
)

/*

Matrix ID_1        -----------   -----     ID_2        -----------   -----     ID_3        -----------   -----
        State_1     -----------   -----     State_2     -----------   -----     State_3     -----------   -----
        FLoor_1     -----------   -----     FLoor_2     -----------   -----     FLoor_3     -----------   -----
        Dir_1       -----------   -----     Dir_2       -----------   -----     Dir_3       -----------   -----
        --------    Down_floor_4  Cab_4     --------    Down_floor_4  Cab_4     --------    Down_floor_4  Cab_4
        Up_floor_3  Down_floor 3  Cab_3     Up_floor_3  Down_floor 3  Cab_3     Up_floor_3  Down_floor 3  Cab_3
        Up_floor_2  Down_floor_2  Cab_2     Up_floor_2  Down_floor_2  Cab_2     Up_floor_2  Down_floor_2  Cab_2
        Up_floor_1  --------      Cab_1     Up_floor_1  --------      Cab_1     Up_floor_1  --------      Cab_1
*/

const numFloors = 4
const numElevators = 3

//ElevatorMatrix := initialize.InitializeMatrix(numFloors,numElevators)  // Set up matrix, add ID

func AddOrder(elevID int, matrix [][]int, buttonPressed io.ButtonEvent) [][]int{
  matrix[7-buttonPressed.Floor][elevID*numFloors + int(buttonPressed.Button)] = 1
  return matrix

}

/*
func isOrderAbove(elevator_floor) {
  for floor := numFloors-; floor< 2; floor++ {

  }
}
*/
