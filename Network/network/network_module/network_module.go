package network_module

import (
	"../bcast"
	"../localip"
	"../peers"
	"flag"
	"fmt"
	"os"
	"time"
)

//defining struct to be transmitted
type TransmitStruct struct {
	Message [3]int
}

func NetworkMod() {
	/*  switch v := Msg.message.(type) {
	    case int:

	    case string:

	    }*/
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)

	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	messageTx := make(chan TransmitStruct)
	messageRx := make(chan TransmitStruct)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, messageTx)
	go bcast.Receiver(16569, messageRx)

	// The example message. We just send one of these every second.
	go func() {
		Msg := TransmitStruct{[3]int{1, 2, 3}}
		for {
			messageTx <- Msg
			time.Sleep(1 * time.Second)
			//			fmt.Println(peerUpdateCh)
			fmt.Println(Msg.Message[2])
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-messageRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}
