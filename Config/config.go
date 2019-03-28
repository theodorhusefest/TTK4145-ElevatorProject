package config

const (
	NumFloors    = 4
	NumElevators = 3
)

type Elevator struct {
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
	IDLE ElevState = iota
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
	Done bool
	SenderID int

	// Select = 1: New order
	// Select = 2: Order complete
	ID     int
	Floor  int
	Button ButtonType

	// Select = 3: Update state/floor/diretion
	State int
	Dir   MotorDirection

	// Select = 4: Acknowledge msg
	Ack bool

	// Select = 5: Resend whole matrix
	ResendMatrix bool
	Matrix       [][]int
}

type AckStruct struct {
	Data        int
	AwaitingAck [NumElevators]bool
}

