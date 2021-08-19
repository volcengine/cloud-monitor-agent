package monitor

import (
	"github.com/volcengine/cloud-monitor-agent/monitor/service"
)

// Service is one of the services of agent.
type Service struct {
}

// Init Service config.
func (s *Service) Init() {
}

// Start collect task.
func (s *Service) Start() {
	go service.CollectToServer()
}
