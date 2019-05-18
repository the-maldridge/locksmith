package main

import (
	"log"
	"net/rpc"

	"github.com/the-maldridge/locksmith/internal/keyhole"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := "wg0"
	reply := keyhole.InterfaceInfo{}

	if err := client.Call("Keyhole.DeviceInfo", args, &reply); err != nil {
		log.Fatal(err)
	}

	log.Println(reply)
}
