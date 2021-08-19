package manager

import (
	"github.com/volcengine/cloud-monitor-agent/monitor"
)

// ServiceManager store the services.
type ServiceManager struct {
	services map[string]Service
}

// NewServiceManager return a service Manager.
func NewServiceManager() *ServiceManager {
	servicesMap := make(map[string]Service)
	return &ServiceManager{services: servicesMap}
}

// RegisterService .
func (sm *ServiceManager) RegisterService() {
	sm.services["monitorService"] = &monitor.Service{}
}

// InitService .
func (sm *ServiceManager) InitService() {
	for _, service := range sm.services {
		service.Init()
	}
}

// StartService .
func (sm *ServiceManager) StartService() {
	for _, service := range sm.services {
		service.Start()
	}
}
