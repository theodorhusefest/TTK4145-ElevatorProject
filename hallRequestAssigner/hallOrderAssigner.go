
package hallOrderAssigner

import (
    "fmt"
    "encoding/json"
    "os/exec"
    . "../Config"
    "strconv"
    )

type HallAssignerElev struct {
    Behaviour    string                             `json:"behaviour"`
    Floor        int                                `json:"floor"`
    Direction    string                             `json:"direction"`
    CabRequests  [NumFloors]bool                    `json:"cabRequests"`
}

type HallAssignerInput struct {
    HallRequests [NumFloors][2]bool                `json:"hallRequests"`       
    States       map[string]*HallAssignerElev      `json:"states"`
}


func AssignHallOrder(newGlobalOrder ButtonEvent, elevatorMatrix [][]int){

    OrderInput := HallAssignerInput{}
    OrderInput.States = make(map[string]*HallAssignerElev)

    // Find all active orders in matrix
    for floor := 0; floor < NumFloors; floor++ {
        for button := 0; button < 2; button ++ {
            for elev := 0; elev < NumElevators; elev++ {
                if elevatorMatrix[len(elevatorMatrix)-floor-1][button+elev*NumElevators] == 1 {
                    OrderInput.HallRequests[floor][button] = true
                } else {
                    OrderInput.HallRequests[floor][button] = false
                }
            }
        }
    }
    OrderInput.HallRequests[newGlobalOrder.Floor][int(newGlobalOrder.Button)] = true

    // Update states
    for elev := 0; elev < NumElevators; elev++ {
        if elevatorMatrix[1][elev*NumElevators] == 3 { // Elevator is offline
            continue
        } else {
            onlineElev := HallAssignerElev{}
            onlineElev.Floor = elevatorMatrix[2][elev*NumElevators]

            if elevatorMatrix[3][elev*NumElevators] == 0 {
                onlineElev.Direction = "stop"
            } else if elevatorMatrix[3][elev*NumElevators] == 1 {
                onlineElev.Direction = "up"
            } else {
                onlineElev.Direction = "down"
            }
            
            // Check cabrequests
            var CabRequests [NumFloors]bool
            for floor := 0; floor < NumFloors; floor++ {
                if elevatorMatrix[len(elevatorMatrix)-floor-1][2+elev*NumElevators] == 1 {
                    CabRequests[floor] = true
                } else {
                    CabRequests[floor] = false
                }
            }
            onlineElev.CabRequests = CabRequests

            if elevatorMatrix[1][elev*NumElevators] == 0 {
                onlineElev.Behaviour = "idle"
            } else if elevatorMatrix[1][elev*NumElevators] == 1 {
                onlineElev.Behaviour = "moving"
            } else {
                onlineElev.Behaviour = "doorOpen"
            }
            ID_str := strconv.Itoa(elevatorMatrix[0][elev*NumElevators])

            OrderInput.States[ID_str] = &onlineElev
        }
    }

    
    fmt.Println("Before encoding: ", OrderInput)

    fmt.Println()

    arg, _ := json.Marshal(OrderInput)
    fmt.Println("After encoding: ", string(arg))

    result, err := exec.Command("sh", "+x" ,"-c", "./hallAssigner -i'"+string(arg)+"'").Output()
    if err != nil {
        fmt.Println("Error in hall Request assigner", err)
    }
    fmt.Println()
    fmt.Println("Result: ", string(result))
    
}
