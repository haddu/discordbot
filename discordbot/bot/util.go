package bot

import (
	"os"
	"os/signal"
	"syscall"
)

func waitForInterupt() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
