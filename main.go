package main

import (
  "fmt"
  "./Initialize"
  "./utilities"
=======
  "fmt"
>>>>>>> origin/FSM
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
=======
    for {
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            io.SetButtonLamp(a.Button, a.Floor, true)

        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            if a == numFloors-1 {
                d = io.MD_Down
            } else if a == 0 {
                d = io.MD_Up
            }
            io.SetMotorDirection(d)


        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                io.SetMotorDirection(io.MD_Stop)
            } else {
                io.SetMotorDirection(d)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := io.ButtonType(0); b < 3; b++ {
                    io.SetButtonLamp(b, f, false)
                }
            }
            io.SetMotorDirection(io.MD_Stop)
        }
    }
}

>>>>>>> origin/FSM
