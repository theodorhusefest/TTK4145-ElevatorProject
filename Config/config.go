package config

const (
  NumFloors = 4
  NumElevators = 3
  ElevID = 0
)

type ElevConfig struct{
  NumFloors int
  NumElevators int
  ElevID int
}

type MotorDirection int

const (
	DIR_Up   MotorDirection = 1
	DIR_Down                = -1
	DIR_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}


type ElevState int

const (
  IDLE ElevState = 0
  MOVING ElevState = 1
  DOOROPEN ElevState = 2
)

type Elevator struct{
  State ElevState
  Floor int
  Dir MotorDirection
}

type Message struct {
  // new order [3]int     [floor, button, value]
  // Order complete
  //  Acknowledge

}
