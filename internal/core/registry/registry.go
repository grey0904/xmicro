package registry

type InstanceInfo struct {
	ID       string
	Address  string
	Port     int
	Metadata map[string]string
}

type Registry interface {
	Register() error
	Deregister() error
	Discover() ([]InstanceInfo, error)
}
