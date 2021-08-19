package manager

// Service define a service should have Init and Start function
type Service interface {
	Init()
	Start()
}
