package hallOrderAssigner

import (
	. "../Config"
	"../Utilities"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type HallAssignerElev struct {
	Behaviour   string          `json:"behaviour"`
	Floor       int             `json:"floor"`
	Direction   string          `json:"direction"`
	CabRequests [NumFloors]bool `json:"cabRequests"`
}

type HallAssignerInput struct {
	HallRequests [NumFloors][2]bool           `json:"hallRequests"`
	States       map[string]*HallAssignerElev `json:"states"`
}

func AssignHallOrder(newGlobalOrder ButtonEvent, elevatorMatrix [][]int) []Message {

	OrderInput := HallAssignerInput{}
	OrderInput.States = make(map[string]*HallAssignerElev)
	var updatedOrders []Message
	var hallRequests [NumFloors][2]bool

	fmt.Println("In hallAssigner")
	utilities.PrintMatrix(elevatorMatrix, NumFloors, NumElevators)

	// Find all active orders in matrix
	for floor := 0; floor < NumFloors; floor++ {
		for button := 0; button < 2; button++ {
			for elev := 0; elev < NumElevators; elev++ {
				//fmt.Println(elevatorMatrix[len(elevatorMatrix)-floor-1][button+elev*NumElevators])
				if elevatorMatrix[len(elevatorMatrix)-floor-1][button+elev*NumElevators] == 1 {
					hallRequests[floor][button] = true
				}
			}
		}
	}
	OrderInput.HallRequests = hallRequests

	if newGlobalOrder.Button != BT_Cab {
		OrderInput.HallRequests[newGlobalOrder.Floor][int(newGlobalOrder.Button)] = true
	}

	// Update states
	for elev := 0; elev < NumElevators; elev++ {
		if elevatorMatrix[1][elev*NumElevators] != 3 { // Elevator is online
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
			var cabRequests [NumFloors]bool
			for floor := 0; floor < NumFloors; floor++ {
				if elevatorMatrix[len(elevatorMatrix)-floor-1][2+elev*NumElevators] == 1 {
					cabRequests[floor] = true
				} else {
					cabRequests[floor] = false
				}
			}
			onlineElev.CabRequests = cabRequests

			if elevatorMatrix[1][elev*NumElevators] == 0 {
				onlineElev.Behaviour = "idle"
			} else if elevatorMatrix[1][elev*NumElevators] == 1 {
				onlineElev.Behaviour = "moving"
			} else {
				onlineElev.Behaviour = "doorOpen"
			}
			IDstr := strconv.Itoa(elevatorMatrix[0][elev*NumElevators])

			OrderInput.States[IDstr] = &onlineElev
		}
	}

	arg, _ := json.Marshal(OrderInput)
	result, err := exec.Command("sh", "+x", "-c", "./MacHallAssigner -i'"+string(arg)+"'").Output()
	if err != nil {
		fmt.Println("Error in Hall Request Assigner", err)
	} else {
		var assignedOrders map[string][][]bool
		json.Unmarshal(result, &assignedOrders)

		for ElevID, orders := range assignedOrders {
			ElevIDint, _ := strconv.Atoi(ElevID)
			for floor := 0; floor < NumFloors; floor++ {
				for button := 0; button < 2; button++ {
					if orders[floor][button] == true && elevatorMatrix[len(elevatorMatrix)-floor-1][button+ElevIDint*NumElevators] == 0 {
						newOrder := Message{Select: 1, ID: ElevIDint, Floor: floor, Button: ButtonType(button)}
						updatedOrders = append(updatedOrders, newOrder)
					}
				}
			}
		}

		fmt.Println("Assigned Orders: ", assignedOrders)
	}
	fmt.Println("Updated Orders: ", updatedOrders)
	return updatedOrders

}
