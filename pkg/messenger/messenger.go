package messenger

type Messenger interface {
	Send(message string) error
	String() string
}
