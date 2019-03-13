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
  BroadcastTicker chan bool
}

func SyncElevator(syncChans SyncElevatorChannels){
//  broadcastTicker(syncChans)

  ticker := time.NewTicker(50 * time.Millisecond)

  for{
    select {
    case <- ticker.C:
      fmt.Println("hei")
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
