package utilities

import (
  "fmt"
  ."../Config"
)

//a function which prints out any given matrix to terminal
func PrintMatrix(matrix [][]int ,n int, m int){

  for i:=0;i<len(matrix);i++{
    fmt.Print(" || ")
    for j:=0; j<len(matrix[0]);j++{
      fmt.Print(matrix[i][j])
      fmt.Print(" ")
      if((j+1)%3 == 0){
        fmt.Print("|| ")
      }
    }
    fmt.Print("\n")
  }
  fmt.Println(" ")
}

func PrintAckMatrix(matrix [NumFloors+4][3*NumElevators]AckStruct, n int, m int){

    for i:=0;i<len(matrix);i++{
      fmt.Print(" || ")
      for j:=0; j<len(matrix[0]);j++{
        fmt.Print(matrix[i][j])
        fmt.Print(" ")
        if((j+1)%3 == 0){
          fmt.Print("|| ")
        }
      }
      fmt.Print("\n")
    }
    fmt.Println(" ")
}

/*
func StateToInt(elevator config.Elevator) int{
  state := elevator.State
	switch state{
	case state.IDLE:
		return 0
	case MOVING:
		return 1
	case DOOROPEN:
		return 2
	}
	return -1
}
*/
