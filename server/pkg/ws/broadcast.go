package ws

import (
	"fmt"
	"slices"
)

var globalBroadcast = make(chan Data)

func GlobalBroadcastHandler() {
	for {
		// get data from channel
		data := <-globalBroadcast
		fmt.Printf("Send Data '%s' to all\n", data.Action)

		// send to all clients
		for _, client := range Users {
			if slices.Contains(data.Excepted, client) {
				fmt.Printf("Excepted \"%s\"\n", client.Username)
				continue
			}
			client.SendRaw(data)
		}
	}
}
