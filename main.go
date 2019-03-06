package main

import (
  "fmt"
  "./Initialize"
  "./utilities"
  "./IO"
)




func main() {

  const numFloors = 4
  const numElevators = 3

  io.Init("localhost:15657",4)

  channelFloor := make(chan int) //channel that is used in InitElevator. Should maybe have a struct with channels?

  elevatorMatrix := initialize.InitializeMatrix(numFloors,numElevators)  // Set up matrix, add ID
  initialize.InitElevator(0,elevatorMatrix,channelFloor)  // Move elevator to nearest floor and update matrix

  utilities.PrintMatrix(elevatorMatrix,numFloors,numElevators)

  floorChn :=  make(chan int)
  go io.PollFloorSensor(floorChn)
  floor := <- floorChn;
  fmt.Println(floor)

}
