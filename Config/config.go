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

type Message struct {
	Select MessageType
	Done bool
	SenderID int

	ID     int
	Floor  int
	Button ButtonType

	State int
	Dir   MotorDirection

	Ack bool

	ResendMatrix bool
	Matrix       [][]int
}

type AckStruct struct {
	Data        int
	AwaitingAck [NumElevators]bool
}

