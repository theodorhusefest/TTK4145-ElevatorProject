package hallOrderAssigner

import (
	. "../Config"
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

func AssignHallOrder(newGlobalOrder ButtonEvent, elevatorMatrix [][]int) {

	OrderInput := HallAssignerInput{}
	OrderInput.States = make(map[string]*HallAssignerElev)

	// Find all active orders in matrix
	for floor := 0; floor < NumFloors; floor++ {
		for button := 0; button < 2; button++ {
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
		if elevatorMatrix[1][elev*NumElevators] != 3 { // Elevator is offline
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

	arg, _ := json.Marshal(OrderInput)
	result, err := exec.Command("sh", "+x", "-c", "./hallAssigner -i'"+string(arg)+"'").Output()
	if err != nil {
		fmt.Println("Error in hall Request assigner", err)
	}

	var assignedOrders map[string][][]bool
	json.Unmarshal(result, &assignedOrders)
	fmt.Println(ButtonType(1))
	var updatedOrders []Message

	for floor := 0; floor < NumFloors; floor++ {
		for button := 0; button < 2; button++ {
			for elev := 0; elev < NumElevators; elev++ {
				//fmt.Println("Floor: ", floor, "button: ", button, assignedOrders[strconv.Itoa(elev)][floor][button])
				if assignedOrders[strconv.Itoa(elev)][floor][button] == true && elevatorMatrix[len(elevatorMatrix)-floor-1][button+elev*NumElevators] == 0 {
					update := Message{ID: elev, Floor: floor, Button: ButtonType(button)}
					updatedOrders = append(updatedOrders, update)
				}
			}
		}
	}
	return updatedOrders

}