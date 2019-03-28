package initialize

import (
	. "../Config"
	"../IO"
	"fmt"
)

//function which puts the posistion of elevator ID into the elevator matrix.
//if elevator is not on a floor, the elevator moves downwards until it hits a sensor
func InitElevator(elevator Elevator, matrix [][]int, channelFloor chan int) {

	io.SetMotorDirection(DIR_Down) //elevator goes downwards
	go io.PollFloorSensor(channelFloor)
	currentFloor := <-channelFloor          //the floor is put onto channelFloor
	matrix[2][elevator.ID*3] = currentFloor //channelFloor is stored in matrix
	elevator.Floor = currentFloor
	io.SetMotorDirection(DIR_Stop) //elevator stops
	InitLights()

}

//creates a matrix with dimensions floor and elevators
func InitializeMatrix() [][]int {
	matrix := make([][]int, 4+NumFloors)
	for i := 0; i < 4+NumFloors; i++ {
		matrix[i] = make([]int, 3*NumElevators)
	}
	for i := 0; i < NumElevators; i++ {
		matrix[1][3*i] = 3 // Set all elevators to offline
	}
	return matrix
}

func InitLights() {
	io.SetDoorOpenLamp(false)
	for floor := 0; floor < NumFloors; floor++ {
		for button := 0; button < 3; button++ {
			io.SetButtonLamp(ButtonType(button), floor, false)
		}
	}
}

func Initialize(numFloors int, numElevators int) ([][]int, Elevator) {

	var inp int
	fmt.Print("Enter elevator ID: ")
	_, err := fmt.Scanf("%d", &inp)
	if err != nil {
		fmt.Println("Error in input!")
	}

	elev := Elevator{ID: inp, State: IDLE}

	elevatorMatrix := InitializeMatrix()
	elevatorMatrix[0][3*elev.ID] = elev.ID
	elevatorMatrix[1][3*elev.ID] = int(IDLE)
	return elevatorMatrix, elev
}

func AssignIDs(matrix [][]int) {
	//assign ID's to elevators
	id := 0
	for i := 0; i < len(matrix[0]); i += 3 {
		matrix[0][i] = id
		id++
	}
}
