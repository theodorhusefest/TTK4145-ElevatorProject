
package main

import "fmt"
import "encoding/json"
import "os"

func main(){
    /*
    msg := {
    "hallRequests" :
        [[false,false],[true,false],[false,false],[false,true]],
    "states" : {
        "one" : {
            "behaviour":"moving",
            "floor":2,
            "direction":"up",
            "cabRequests":[false,false,true,true]
        },
        "two" : {
            "behaviour":"idle",
            "floor":0,
            "direction":"stop",
            "cabRequests":[false,false,false,false]
            }
        }
    }
    */

    type ElevatorStruct struct{ 
        Behaviour string
        Floor int
        Direction string
        CabRequests []bool
    }

    type Elevators struct{
        ElevOne ElevatorStruct
        ElevTwo ElevatorStruct 
    }
    type Message struct{
        HallRequests [][]bool
        States Elevators
    }

    One := ElevatorStruct{
        Behaviour: "moving",
        Floor: 2,
        Direction: "up",
        CabRequests: []bool{false,false,true,true},
    }
    Two := ElevatorStruct{
        Behaviour: "idle",
        Floor: 0,
        Direction: "stop",
        CabRequests: []bool{false, false, false, false},
    }

    ElevStates := Elevators{
        ElevOne: One,
        ElevTwo: Two,
    }

    Msg := Message{
        HallRequests: [][]bool{{false, false},{true, false},{false, false},{false, true}},
        States: ElevStates,
    }
    fmt.Println(Msg)
    fmt.Println()

    
    jsonMsg, err := json.Marshal(Msg)
    if err != nil {
        fmt.Println("Error")
    }
    os.Stdout.Write(jsonMsg)

}

