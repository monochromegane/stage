package stage

import "os"

type Stage struct {
	outDir string
}

func New(outDir string) *Stage {
	return &Stage{
		outDir: outDir,
	}
}

func (s *Stage) Run(iter int, newActorFn NewActorFn, scenario Scenario) error {
	err := ensureOutDir(outDir)
	if err != nil {
		return err
	}

	for i := 0; i < iter; i++ {
		actor := NewActorFn()
		s.run(i, actor, scenario)
	}
	return nil
}

func (s *Stage) run(iter int, actor Actor, scenario Scenario) {
	for scenario.Scan() {
		actor.Act(scenario.Line())
	}
}

func ensureOutDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
