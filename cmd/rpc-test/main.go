package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/earendil-works/pi-mc/internal/rpc"
)

func main() {
	// We use "go run ./cmd/mock-pi" to simulate the pi executable
	client, err := rpc.NewPiRpcClient("go", "run", "./cmd/mock-pi")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	
	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
	defer client.Stop()
	
	fmt.Println("RPC Client started. Sending prompt in 1 second...")
	time.Sleep(1 * time.Second)
	
	err = client.SendCommand(rpc.Command{
		Type: "prompt",
		Payload: map[string]string{"text": "Hello pi!"},
	})
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}
	
	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	
	for {
		select {
		case ev, ok := <-client.Events():
			if !ok {
				fmt.Println("Event channel closed")
				return
			}
			if ev.Type == "message_update" {
				fmt.Printf("Received message_update payload: %s\n", string(ev.Payload))
			} else {
				fmt.Printf("Received event: %s\n", ev.Type)
			}
		case err := <-client.Errors():
			fmt.Printf("Error: %v\n", err)
		case <-sigChan:
			fmt.Println("\nExiting...")
			return
		}
	}
}
