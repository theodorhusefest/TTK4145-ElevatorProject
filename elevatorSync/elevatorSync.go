package syncElevator

import (
  "fmt"
  //"time"
  "../Network/network/peers"
  "../Config"
)



type SyncElevatorChannels struct{
  OutGoingOrder chan config.Message
  MessageToSend chan config.Message
  InComingOrder chan config.Message
  MessageRecieved chan config.Message
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
    //Transmit message from orderManager
    case msg := <- syncChans.MessageToSend:
      message.ID = msg.ID
      message.Floor = msg.Floor
      message.Button = msg.Button
      syncChans.OutGoingOrder <- message

    case msg := <- syncChans.InComingOrder:
      message.ID = msg.ID
      message.Floor = msg.Floor
      message.Button = msg.Button
      syncChans.MessageRecieved <- message


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
