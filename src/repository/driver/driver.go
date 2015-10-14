package driver

type Driver interface {
	Run() error
	GetName() string
	GetSource() map[string]string
	GetReference() string
}
