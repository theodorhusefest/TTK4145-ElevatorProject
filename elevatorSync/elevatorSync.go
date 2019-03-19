package syncElevator

import (
  "fmt"
  "time"
  "../Network/network/peers"
  "../Config"
)



type SyncElevatorChannels struct{
  OutGoingMsg chan config.Message
  InCommingMsg chan config.Message
  ChangeInOrderch chan config.Message
  PeerUpdate chan peers.PeerUpdate
  TransmitEnable chan bool
  BroadcastTicker chan bool
}

func SyncElevator(syncChans SyncElevatorChannels, elevatorConfig config.ElevConfig, UpdateElevatorChan chan config.Message){

  var online bool

//  broadcastTicker(syncChans)
  message := config.Message{
  }

  broadCastTicker := time.NewTicker(100 * time.Millisecond)

  for{
    select {
    // --------------------------------------------------------------------------Case triggered by broadcast-ticker.
    case <- broadCastTicker.C:
      if online {
        syncChans.OutGoingMsg <- message
//        fmt.Println(message)
      }



    // --------------------------------------------------------------------------Case triggered by local ordermanager, change in order
    case changeInOrder := <-syncChans.ChangeInOrderch:
      //Håndter endring som kom fra ordermanager. Send alt inn på message og sett message.Done = false
      message = changeInOrder

      // Broadcast message
      syncChans.OutGoingMsg <- message

      // Vent til alle er enige, gi klarsignal til ordermanager ??????

      // Sett message.Done = true
      message.Done = true



    // --------------------------------------------------------------------------Case triggered by bcast.Recieving
    case msgRecieved := <- syncChans.InCommingMsg:

      // Check if message has been processed.
      if !msgRecieved.Done {
        message = msgRecieved

        // If select = 1, new order was recieved.
        if message.Select == 1 {
          // Sett info inn på message
        }

      // Wait to everyone agree

      // Send message to local ordermanager
      UpdateElevatorChan <- message
      message.Done = true
      }



    // --------------------------------------------------------------------------Case triggered by update in peers
    case peer := <- syncChans.PeerUpdate:
    //Update peers
    //Check how many peers are connected
    //If only you, start singelmode ???????????????????

      if (len(peer.Peers) == 0) {
        fmt.Println("I'm offline")
        online = false
      } else {
        fmt.Println("I'm online")
        online = true
      }
      fmt.Println(peer.Peers)

      

      // ????????????    orderManager.addOrder(elevatorConfig,peer.Peers[string(elevatorConfig.ElevID)],)     ?????????????????????


    }
  }
}
