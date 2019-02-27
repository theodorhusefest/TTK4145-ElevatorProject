package main

import (
  //"fmt"
  "./Initialize"
  "./utilities"
)




func main() {

  const floors = 4
  const elevators = 5

  elevatorMatrix := initialize.InitializeMatrix(floors,elevators)
  initialize.AssignIDs(elevatorMatrix)
  utilities.PrintMatrix(elevatorMatrix,floors,elevators)
}
