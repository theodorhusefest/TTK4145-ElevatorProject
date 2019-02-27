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

  return matrix
}

func AssignIDs(matrix [][]int){
  //assign ID's to elevators
  id := 1
  for i:=0; i<len(matrix[0]);i+=3{
    matrix[0][i] = id
    id++
  }
}

//func AssignStates(matrix [][]int){
