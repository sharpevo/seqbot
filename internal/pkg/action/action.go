package action

type ActionInterface interface {
	Run(eventname string, command CommandInterface) (string, error)
	Name() string
}
