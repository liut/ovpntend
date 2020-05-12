package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/thejerf/suture"
)

// Servicer ...
type Servicer interface {
	Serve()
	Stop()
}

var (
	services []Servicer
)

func addService(svr Servicer) {
	services = append(services, svr)
}

func startUp() {
	count := len(services)
	if count == 0 {
		return
	}
	supervisor := suture.NewSimple("main")

	idleConnsClosed := make(chan struct{})
	go func() {
		logger().Info("waiting signal")
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill, syscall.SIGTERM)
		sig := <-sigint

		logger().Infow("received a signal, shuting down", "sig", sig)
		supervisor.Stop()
		close(idleConnsClosed)
	}()

	for i := 0; i < count; i++ {
		supervisor.Add(services[i])
	}
	supervisor.ServeBackground()

	<-idleConnsClosed
}
