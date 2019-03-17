package initialize

import (
  "fmt"
  //"../utilities"
    //"../FSM"
  "../IO"
  "../Config"


)



//creates a matrix with dimensions floor and elevators
func InitializeMatrix( floors int, elevators int)[][]int{
  matrix := make([][]int,4+floors)
  for i:=0; i< 4+floors;i++{
    matrix[i] = make([]int,3*elevators)
  }

  return matrix
}


func Initialize() ([][]int, config.ElevConfig){
  var inp int
  conf := config.ElevConfig{NumFloors:0,NumElevators:0,ElevID:0}
	fmt.Print("Enter NumElevators: ")
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

	fmt.Print("Enter elevator ID: ")
  _,err=fmt.Scanf("%d", &inp)
  if(err!=nil){
    fmt.Println("Error in input!")
  }
  conf.ElevID = inp

  elevatorMatrix := InitializeMatrix(conf.NumFloors,conf.NumElevators)
  elevatorMatrix[0][3*conf.ElevID] = conf.ElevID
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
func InitElevator(conf config.ElevConfig, matrix [][]int, channelFloor chan int){

  io.SetMotorDirection(-1) //elevator goes downwards
  go io.PollFloorSensor(channelFloor) //the floor is put onto channelFloor
  matrix[2][conf.ElevID*3] = <-channelFloor //channelFloor is stored in matrix
  io.SetMotorDirection(0) //elevator stops
}
