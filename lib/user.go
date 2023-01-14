package lib

type User interface {
	GetID() string
	GetExtra() map[string]string
}
