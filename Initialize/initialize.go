package initialize

import (
  "fmt"
  "../IO"
  . "../Config"


)



//creates a matrix with dimensions floor and elevators
func InitializeMatrix() [][]int {
  matrix := make([][]int, 4+NumFloors)
  for i:=0; i < 4+NumFloors; i++{
    matrix[i] = make([]int, 3*NumElevators)
  }
  for i := 0; i < NumElevators; i++ {
    matrix[1][3*i] = 3 // Set all elevators to offline
  }
  return matrix
}


func Initialize(numFloors int, numElevators int) ([][]int, ElevConfig){


  /*fmt.Print("Enter NumElevators: ")
	_,err := fmt.Scanf("%d", &inp)
  if(err!=nil){
    fmt.Println("Error in input!")
  }
  conf.NumElevators = inp

	fmt.Print("Enter NumFloors: ")
  _,err=fmt.Scanf("%d", &inp)
  if(err!=nil){
    fmt.Println("Error in input!")
  }
  conf.NumFloors = inp
*/
  var inp int
	fmt.Print("Enter elevator ID: ")
  _ , err := fmt.Scanf("%d", &inp)
  if (err != nil) {
    fmt.Println("Error in input!")
  }

  //conf.ElevID = inp
  conf := ElevConfig{NumFloors:numFloors, NumElevators:numElevators, ElevID:inp,}

  elevatorMatrix := InitializeMatrix()
  elevatorMatrix[0][3*conf.ElevID] = conf.ElevID
  elevatorMatrix[1][3*conf.ElevID] = 0                // Set state to Idle
  return elevatorMatrix, conf
}



func AssignIDs(matrix [][]int){
  //assign ID's to elevators
  id := 0
  for i:=0; i<len(matrix[0]);i+=3{
    matrix[0][i] = id
    id++
  }
}

//function which puts the posistion of elevator ID into the elevator matrix.
//if elevator is not on a floor, the elevator moves downwards until it hits a sensor
func InitElevator(conf ElevConfig, matrix [][]int, channelFloor chan int){

  io.SetMotorDirection(-1) //elevator goes downwards
  go io.PollFloorSensor(channelFloor) //the floor is put onto channelFloor
  matrix[2][conf.ElevID*3] = <-channelFloor //channelFloor is stored in matrix
  io.SetMotorDirection(0) //elevator stops
}
