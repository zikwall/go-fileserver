package main

import (
	"github.com/zikwall/go-fileserver/src/lib"
	"os"
	"os/signal"
	"syscall"
)

func congratulations() {
	lib.Info("Congratulations, the file server has been successfully launched")
}

func waitSystemNotify() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig
}
