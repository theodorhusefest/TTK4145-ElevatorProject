package main

import (
  //"fmt"
  "./Initialize"
  "./utilities"
  "./IO"
)




func main() {

  const floors = 4
  const elevators = 3

  io.Init("localhost:15600",4)

  channelFloor := make(chan int) //channel that is used in InitElevator. Should maybe have a struct with channels?



  elevatorMatrix := initialize.InitializeMatrix(floors,elevators)
  initialize.AssignIDs(elevatorMatrix)
  initialize.InitElevator(0,elevatorMatrix,channelFloor)

  utilities.PrintMatrix(elevatorMatrix,floors,elevators)
}
