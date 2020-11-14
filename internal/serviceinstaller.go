package internal

import (
	"fmt"
	"github.com/kardianos/service"
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

func InstallServiceIfNotYetInstalled(c *Configuration) bool {
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
		fmt.Println(err)
	}

	err = s.Install()
	if err != nil {
		fmt.Println(err)
	}
	err = s.Start()
	if err != nil {
		fmt.Println(err)
	}

	return true
}