
package main

import "fmt"
import "encoding/json"
import "os/exec"
//import "os"
//import "path/filepath"
//import "sync"

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
    //fmt.Println(Msg)
    fmt.Println()
    
    jsonMsg, err := json.Marshal(Msg)
    if err != nil {
        fmt.Println("Error")
    }
    fmt.Println()
    */

    type AssignerCompatibleElev struct {
        //sync.RWMutex `json:"-"`
        Behaviour    string  `json:"behaviour"`
        Floor        int     `json:"floor"`
        Direction    string  `json:"direction"`
        CabRequests  [4]bool `json:"cabRequests"`
    }

    type AssignerCompatibleInput struct {
       // sync.RWMutex `json:"-"`
        HallRequests [4][2]bool                         `json:"hallRequests"`
        States       map[string]*AssignerCompatibleElev `json:"states"`
    }

    Elev1 := AssignerCompatibleElev{
        Behaviour: "moving",
        Floor: 2,
        Direction: "up",
        CabRequests: [4]bool{false,false,true,true},
    } 

    Elev2 := AssignerCompatibleElev{
        Behaviour: "idle",
        Floor: 0,
        Direction: "stop",
        CabRequests: [4]bool{false, false, false, false},
    }

    Elevs := AssignerCompatibleInput{}
    Elevs.HallRequests = [4][2]bool{{false, false},{true, false},{false, false},{false, true}}
    Elevs.States = make(map[string]*AssignerCompatibleElev)
    Elevs.States["one"] = &Elev1
    Elevs.States["two"] = &Elev2
    
    fmt.Println("Before encoding: ", Elevs)

    fmt.Println()

    arg, _ := json.Marshal(Elevs)
    fmt.Println("After encoding: ", string(arg))

    result, err := exec.Command("sh", "+x" ,"-c", "./hallRequestAssigner -i'"+string(arg)+"'").Output()
    if err != nil {
        fmt.Println("Error in hall Request assigner", err)
    }

    //var a map[string][][]bool
    //json.Unmarshal(result, &a)
    fmt.Println()
    fmt.Println("Result: ", string(result))
    
}
