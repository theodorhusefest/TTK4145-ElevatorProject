package syncElevator

import (
  "fmt"
  "time"
  "strconv"
  "../Network/network/peers"
  . "../Config"

)



type SyncElevatorChannels struct{
  OutGoingMsg chan []Message
  InCommingMsg chan []Message
  ChangeInOrderch chan []Message
  SendFullMatrixch chan [][]int
  PeerUpdate chan peers.PeerUpdate
  TransmitEnable chan bool
  BroadcastTicker chan bool
  AskForMatrix chan bool
}

func SyncElevator(syncChans SyncElevatorChannels, elevatorConfig ElevConfig, UpdateElevatorChan chan []Message){

  //var online bool

//  broadcastTicker(syncChans)

  broadCastTicker := time.NewTicker(100 * time.Millisecond)
  online := false
  for{
    select {
    // --------------------------------------------------------------------------Case triggered by broadcast-ticker.
    /*case <- broadCastTicker.C:
      if online {
        syncChans.OutGoingMsg <- message
//        fmt.Println(message)
      }

*/

    // --------------------------------------------------------------------------Case triggered by local ordermanager, change in order
    case changeInOrder := <- syncChans.ChangeInOrderch:
      //Håndter endring som kom fra ordermanager. Send alt inn på message og sett message.Done = false
      //message := changeInOrder


      // Broadcast message
      select {
      case <- broadCastTicker.C:
          syncChans.OutGoingMsg <- changeInOrder
          fmt.Println("Sending message")
      }


      // Vent til alle er enige, gi klarsignal til ordermanager ??????

      // Sett message.Done = true



    // --------------------------------------------------------------------------Case triggered by bcast.Recieving
    case msgRecieved := <- syncChans.InCommingMsg:
        /*
      // Check if message has been processed.
      if !msgRecieved.Done {
        message := msgRecieved

        // If select = 1, new order was recieved.
        if message.Select == 1 {
          // Sett info inn på message
        }

      // Wait to everyone agree

      // Send message to local ordermanager*/
      UpdateElevatorChan <- msgRecieved





      //}

    // --------------------------------------------------------------------------Case triggered by update in peers
    case p := <- syncChans.PeerUpdate:
    if len(p.Peers) == 0 {
        online = false
    }



    //Update peers
    //Check how many peers are connected
    //If only you, start singelmode ???????????????????
    fmt.Println("Peers: ", p.Peers)
    for _, peersOnline := range p.Peers {
        newID, _ := strconv.Atoi(peersOnline)
        if (elevatorConfig.OnlineList[newID] == false) && !online {
            fmt.Println(peersOnline)
            fmt.Println(elevatorConfig.OnlineList[newID])
            elevatorConfig.OnlineList[newID] = true
            fmt.Println("Ask for resend Matrix")
            message := []Message {{Select: 5, ID: newID}}
            //syncChans.OutGoingMsg <- message
            fmt.Println(message)
            online = true
        } else if (elevatorConfig.OnlineList[newID] == false) {
            elevatorConfig.OnlineList[newID] = true
        }
    }
    for _, peersOffline := range p.Lost {
        newID, _ := strconv.Atoi(peersOffline)
        elevatorConfig.OnlineList[newID] = false
    }

            //Send all matrix
            // if has not been online before, add to online list

    fmt.Println(elevatorConfig.OnlineList)

    /*
      if (len(peer.Peers) == 0) {
        fmt.Println("I'm offline")
        //online = false
      } else {
        fmt.Println("I'm online")
        fmt.Println("Currently online:", peer.Peers)
        fmt.Println(peer.New, "just connected")

        //online = true
      }
*/



      // ????????????    orderManager.addOrder(elevatorConfig,peer.Peers[string(elevatorElevID)],)     ?????????????????????


    }
  }
}
