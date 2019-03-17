package syncElevator

import (
  "fmt"
  //"time"
  "../Network/network/peers"
  "../Config"
)



type SyncElevatorChannels struct{
  OutGoingMsg chan config.Message
  InCommingMsg chan config.Message
  ChangeInOrder chan config.Message
  PeerUpdate chan peers.PeerUpdate
  TransmitEnable chan bool
  BroadcastTicker chan bool
}

func SyncElevator(syncChans SyncElevatorChannels, elevatorConfig config.ElevConfig){
//  broadcastTicker(syncChans)
  message := config.Message{
  }

//  broadCastTicker := time.NewTicker(500 * time.Millisecond)

  for{
    select {
    case ChangeInOrder := <- syncChans.InCommingMsg:

    // Check if incomming message is of order-done. If so, remove that order
    if ChangeInOrder.Done {
      message.Floor = 1
    } else if ChangeInOrder.Select == 1 {
      // ADD ORDER TO LOCAL ELEVATOR VIA NewNetworkOrder channel
    }












    //Transmit message from orderManager


    //Update peers
    case peer := <- syncChans.PeerUpdate:
      fmt.Println(peer.Peers)
      //orderManager.addOrder(elevatorConfig,peer.Peers[string(elevatorConfig.ElevID)],)
    }
/*
    select{

    //New local order, insert into msg for transmitting
    case newOrderRecieved := <-:
      fmt.Println("newlocalorder"+ newLocalOrder)




    //If elevator online, send msg on channel for BCAST-transmitter
    case <-syncChans.broadcastTicker:
      fmt.Println("broadcastticker")






    //BCAST-reciever recieve new message
    case msg := <-NewOrderRecieved:
      fmt.Println("incomingOrder")


    //Check how many peers are connected
    case peer <- PeerUpdate:
    //If only you, start singelmode

    }

*/
  }
}
