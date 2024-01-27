package lib

type MessageSource interface {
	Receive(batch int32) ([]byte, error)
}
