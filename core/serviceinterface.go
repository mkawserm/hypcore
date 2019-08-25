package core

type ServiceInterface interface {
	ServiceName() string
	ServiceGroup() string

	//Set Service Settings
	SetServiceSettings(data []byte)

	// Serve based on data and return boolean flag with proper message
	Serve(data []byte) (bool, []byte)
}
