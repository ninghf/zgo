package service

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zsounder/zgo/logger"
)

type ServiceIF interface {
	OnInit()
	OnClose()
	Run(closeChan chan struct{})
	Name() string
}

func Run(services ...ServiceIF) {
	Registe(services...)
	runAll()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT)

	logger.Debug("Progress ", os.Getpid(), " start")

	// Block until a signal is received
	sig := <-c

	logger.Debug("closing down by signal:", sig)

	closeAll()
}

func Registe(sifs ...ServiceIF) {
	for i := 0; i < len(sifs); i++ {
		registerOne(sifs[i])
	}
}

type serviceAgent struct {
	sif       ServiceIF
	wg        sync.WaitGroup
	closeChan chan struct{}
}

var allAgents []*serviceAgent

func registerOne(sif ServiceIF) {
	allAgents = append(allAgents, &serviceAgent{sif: sif, closeChan: make(chan struct{}, 1)})
}

func runAll() {
	for i := 0; i < len(allAgents); i++ {
		allAgents[i].sif.OnInit()
	}

	for i := 0; i < len(allAgents); i++ {
		go func(a *serviceAgent) {
			runService(a)
		}(allAgents[i])
	}
}

func closeAll() {
	for i := len(allAgents) - 1; i >= 0; i-- {
		agent := allAgents[i]
		agent.closeChan <- struct{}{}
		agent.wg.Wait()
		destroyService(agent)
	}
}

func runService(agent *serviceAgent) {
	agent.wg.Add(1)
	agent.sif.Run(agent.closeChan)
	agent.wg.Done()
}

func destroyService(agent *serviceAgent) {
	agent.sif.OnClose()
}
