package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Listen(onShutdown func()) {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		sig := <-sigCh
		log.Printf("received shutdown signal: %s", sig)
		onShutdown()
	}()
}