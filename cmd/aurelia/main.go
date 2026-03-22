package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "onboard":
			if err := runOnboard(os.Stdin, os.Stdout); err != nil {
				log.Fatalf("Failed to run onboarding: %v", err)
			}
			return
		case "cron":
			if err := runCronCLI(os.Args[2:]); err != nil {
				log.Fatalf("Cron command failed: %v", err)
			}
			return
		case "telegram":
			if err := runTelegramCLI(os.Args[2:]); err != nil {
				log.Fatalf("Telegram command failed: %v", err)
			}
			return
		default:
			log.Fatalf("Unknown command: %s", os.Args[1])
		}
	}

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
