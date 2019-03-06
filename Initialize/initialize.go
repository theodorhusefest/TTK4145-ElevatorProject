package initialize

import (
  //"fmt"
  //"../utilities"
    //"../FSM"
  "../IO"
)



//creates a matrix with dimensions floor and elevators
func InitializeMatrix( floors int, elevators int)[][]int{
  matrix := make([][]int,4+floors)
  for i:=0; i< 4+floors;i++{
    matrix[i] = make([]int,3*elevators)
  }
  AssignIDs(matrix)
  return matrix
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
func InitElevator(elevID int,matrix [][]int, channelFloor chan int){

  io.SetMotorDirection(-1) //elevator goes downwards
  go io.PollFloorSensor(channelFloor) //the floor is put onto channelFloor
  matrix[2][elevID*3] = <-channelFloor //channelFloor is stored in matrix
  io.SetMotorDirection(0) //elevator stops
}
