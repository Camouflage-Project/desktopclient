package internal

import (
	"github.com/kardianos/service"
	"go.uber.org/zap"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	return nil
}
func (p *program) run() {
}
func (p *program) Stop(s service.Service) error {
	return nil
}

func InstallServiceIfNotYetInstalled(c *Configuration,  logger *zap.Logger) bool {
	if !service.Interactive() {
		return false
	}

	svcConfig := &service.Config{
		Name:        c.CurrentVersion,
		DisplayName: c.CurrentVersion,
		Description: c.CurrentVersion,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)

	if err != nil {
		logger.Error(err.Error())
	}else {
		err = s.Install()
		if err != nil {
			logger.Error(err.Error())
		}

		err = s.Start()
		if err != nil {
			logger.Error(err.Error())
		}
	}

	logger.Info("installed desktopClient as service")
	return true
}