package action

type ActionInterface interface {
	Run(wfqLogPath string, chipId string) (string, error)
	Name() string
}
