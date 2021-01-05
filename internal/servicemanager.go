package internal

import (
	"desktopClient/config"
	"github.com/kardianos/service"
	"go.uber.org/zap"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	StartLogic()
	//TestStartLogic()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func InstallServiceIfNotYetInstalled(c *config.Configuration,  logger *zap.Logger) bool {
	if !service.Interactive() {
		return false
	}

	s := initService(c, logger)

	err := s.Install()
	if err != nil {
		logger.Error(err.Error())
	}

	err = s.Start()
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("installed desktopClient as service")
	return true
}

func RunService(c *config.Configuration, logger *zap.Logger) {
	logger.Info("running service")
	s := initService(c, logger)
	err := s.Run()

	if err != nil {
		logger.Error(err.Error())
	}
}

func initService(c *config.Configuration, logger *zap.Logger) service.Service {
	svcConfig := &service.Config{
		Name:        c.CurrentVersion,
		DisplayName: c.CurrentVersion,
		Description: c.CurrentVersion,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)

	if err != nil {
		logger.Error(err.Error())
	}

	return s
}