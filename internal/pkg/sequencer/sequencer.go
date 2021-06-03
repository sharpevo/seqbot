package sequencer

type SequencerInterface interface {
	GetSlide(string) (string, error)
	GetBarcode(string) (string, error)
	GetExtraExperimentInfo(string) (string, error)
	GetArchiveDir(string) (string, error)
	GetResultDir(string) (string, error)
	GetWfqTime(string) (string, error)
	GetUploadTime(string) (string, error)
	IsSuccess(string) (bool, error)
}
