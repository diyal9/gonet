package main

import (
	"gonet/base"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"gonet/message"
)

var(
)

func ExampleDial() {
	origin := "http://localhost/"
	url := "ws://192.168.215.107:31700/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	AccountName := fmt.Sprintf("test%d", 1)
	packet1 := &message.C_A_LoginRequest{PacketHead: message.BuildPacketHead(0, message.SERVICE_ACCOUNTSERVER),
		AccountName: AccountName, BuildNo: base.BUILD_NO}
	buff := message.Encode(packet1)
	buff = base.SetTcpEnd(buff)
	if _, err := ws.Write(buff); err != nil {
		log.Fatal(err)
	}

	for{
		var msg = make([]byte, 512)
		var n int
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received: %s.\n", msg[:n])
	}
}

func main() {
	ExampleDial()
}
