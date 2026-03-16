package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := bootstrapApp()
	if err != nil {
		log.Fatalf("Failed to bootstrap Aurelia: %v", err)
	}
	defer app.close()

	app.start()
	waitForShutdownSignal()

	log.Println("Shutting down Aurelia...")
	app.shutdown(context.Background())
}

func waitForShutdownSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}


