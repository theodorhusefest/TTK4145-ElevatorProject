package syncElevator

import (
  "fmt"
  "time"
  "../Network/network/peers"
//  "../Config"
)



type SyncElevatorChannels struct{
//  OutGoingOrder chan ??
//  InComingOrder chan ??
  PeerUpdate chan peers.PeerUpdate
  TransmitEnable chan bool
  BroadcastTicker chan bool
}

func SyncElevator(syncChans SyncElevatorChannels){
//  broadcastTicker(syncChans)

  broadCastTicker := time.NewTicker(500 * time.Millisecond)

  for{
    select {
    case <- broadCastTicker.C:
      fmt.Println("hei")
    case peer := <- syncChans.PeerUpdate:
      fmt.Println(peer.Peers)
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
