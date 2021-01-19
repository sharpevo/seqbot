package action

type ActionInterface interface {
	Run(eventname string, wfqLogPath string, chipId string) (string, error)
	Name() string
}
