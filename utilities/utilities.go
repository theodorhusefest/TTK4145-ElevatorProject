package utilities

import (
  "fmt"
)


func PrintMatrix(matrix[8][9]int ,n int, m int){
  for i := 0; i<8; i++{
    for j:= 0; j<9; j++{
      fmt.Print(matrix[i][j])
      fmt.Print(" ")
    }
    fmt.Print("\n")
  }
}
