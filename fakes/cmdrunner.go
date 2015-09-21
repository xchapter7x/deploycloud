package fakes

type FakeCmdRunner struct {
	CmdSpy     []string
	RunSpy     int
	ErrRunFake error
}

func (s *FakeCmdRunner) Run() error {
	s.RunSpy++
	return s.ErrRunFake
}
