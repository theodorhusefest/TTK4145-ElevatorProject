package main

import (
  "fmt"
)




func main() {
  elevatorMatrix := [4+NFloors][3*elevators] int{}

  for i:=0; i<4+NFloors; i++ {
     for j := 0; j<3*elevators; j++{
       elevatorMatrix[i][j] = 0
     }
   }
   utilities.PrintMatrix(elevatorMatrix,4+NFloors,elevators)
}
