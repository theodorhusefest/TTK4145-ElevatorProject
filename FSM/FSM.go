package FSM

type state int

const (
  IDLE state = 0
  MOVING state = 1
  DOOROPEN state = 2
)

switch state{
case IDLE:
  //blablabla
case MOVING:
  //blablabla
case DOOROPEN:
  //blablabla
}
