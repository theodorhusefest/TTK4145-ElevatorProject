package config

const (
	NumFloors    = 4
	NumElevators = 3
)

type Elevator struct {
  NumFloors int
  NumElevators int
  ID    int
  State ElevState
  Floor int
  Dir   MotorDirection
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
	IDLE     ElevState = iota
	MOVING  
	DOOROPEN 
	UNDEFINED  
)


type MessageType int

const (
	NewOrder MessageType = iota + 1
	OrderComplete
	UpdateStates
	UpdateOffline
	ACK
	SendMatrix
	UpdatedMatrix
)

// Struct for sendig of messages
type Message struct {
	// Select different cases of message based on value
	Select MessageType
	// Message.Done = true: Message has been processed
	Done bool

	SenderID int

	// Select = 1: NEW ORDER
	// Select = 2: AN ORDER HAS BEEN EXCECUTED
	ID     int
	Floor  int
	Button ButtonType

	// Select = 3: UPDATE STATE/FLOOR/DIR TO A GIVEN ELEVATOR
	State int
	Dir   MotorDirection

	// Select = 4: ACKNOWLEDGE
	Ack bool

	// Select = 5: Ask for whole matrix
	ResendMatrix bool
	Matrix       [][]int

	// Select = 6 Resend message

}


type AckStruct struct {
	Data int
	Data2 int
	AwaitingAck [NumElevators]bool
	RecievedAck [NumElevators]bool
}


