package network
/*
import (
  //"./localip"
  "./peers"
  "./bcast"
  //"./conn"
  "fmt"
)

const PORT = 15602

type NetworkMsg struct{
  message int
}

type NetworkChans struct{
  TxToggleChan chan bool
  //PeersRxToggleChan chan bool
  RxUpdateChan chan int

  RxMessageChan chan int
  TxMessageChan chan int
}

//ID must be implemented

func TransmitMessage(message NetworkMsg, NetChans NetworkChans){
  NetChans.TxMessageChan <-message.message
}

func RecieveMessage(NetChans NetworkChans){
  select{
  case msg := <-NetChans.TxMessageChan:
    fmt.Println(msg)
  }
}
*/
