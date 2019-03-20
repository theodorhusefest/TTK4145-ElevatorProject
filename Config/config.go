package config

const (
  NumFloors = 4
  NumElevators = 3
)

type ElevConfig struct{
  NumFloors int
  NumElevators int
  ElevID int
  OnlineList [NumElevators]bool
  IsOnline bool
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
  OFFLINE ElevState = 3
)

type Elevator struct{
  ID int
  State ElevState
  Floor int
  Dir MotorDirection
}


// Struct for sendig of messages
type Message struct {
  // Select different cases of message based on value
  Select int
  // Message.Done = true: Message has been processed
  Done bool

  // Select = 1: NEW ORDER
  // Select = 2: AN ORDER HAS BEEN EXCECUTED
  ID int
  Floor int
  Button ButtonType

  // Select = 3: UPDATE STATE/FLOOR/DIR TO A GIVEN ELEVATOR
  State int
  Dir MotorDirection

  // Select = 4: ACKNOWLEDGE
  Ack bool

  // Select = 5: Ask for whole matrix
  ResendMatrix bool
  Matrix [][]int

  // Select = 6 Resend message

}
