package utilities

import (
  "fmt"
)

//a function which prints out any given matrix to terminal
func PrintMatrix(matrix [][]int ,n int, m int){

  for i:=0;i<len(matrix);i++{
    for j:=0; j<len(matrix[0]);j++{
      fmt.Print(matrix[i][j])
      fmt.Print(" ")
    }
    fmt.Print("\n")
  }

}
