package FSM

import (
	"../IO"
	"time"

)

type state int

const (
  IDLE state = 0
  MOVING state = 1
  DOOROPEN state = 2
)
func StateMachine(){
  switch state{
  case IDLE:
    //blablabla
  case MOVING:
    //blablabla
  case DOOROPEN:
    //blablabla
  }
}



type FSMChannels struct {
	newOrder 		chan ButtonEvent
	floorReached 	chan int
}


func FSM() {

	doorOpenTimer := newTimer(3 * time.Second)
	doorOpenTimer.Stop()

	for {

		select {
			case 
	

		}
	



	}



}