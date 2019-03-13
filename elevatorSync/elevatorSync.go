package syncElevator

import (
  "fmt"
  "time"
)


type SyncElevatorChannels struct{
  OutGoingOrder chan ??
  InComingOrder chan ??
  PeerUpdate chan ??
}

func SyncElevator(syncChans SyncChannels){
  broadcastTicker(syncChans)


  for{
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
  }
}



func broadcastTicker(syncChans SyncChannels){
  timer := time.NewTimer(5*time.Millisecond)
  <-timer.C
  syncChans.broadcastTicker<-true
}
