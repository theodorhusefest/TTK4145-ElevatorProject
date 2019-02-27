package utilities

import (
  "fmt"
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

}
