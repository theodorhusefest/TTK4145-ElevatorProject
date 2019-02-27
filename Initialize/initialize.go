package main

import (
  "fmt"
  "../utilites/utilites.go"
)

const NFloors = 4
const elevators = 3

elevatorMatrix := [4+NFloors][3*elevators] int{}

func main() {



   for i:=0; i<4+NFloors; i++ {
     for j := 0; j<3*elevators; j++{
       elevatorMatrix[i][j] = 0
     }
   }

   fmt.Println(elevatorMatrix[2][2])

   utilites.printMatrix(elevatorMatrix,4+NFloors,elevators)
}
