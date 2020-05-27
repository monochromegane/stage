package stage

type Stage struct {
}

func New() *Stage {
	return &Stage{}
}

func (s *Stage) Run(iter int, newActorFn NewActorFn, scenario Scenario) {
	for i := 0; i < iter; i++ {
		actor := NewActorFn()
		s.run(i, actor, scenario)
	}
}

func (s *Stage) run(iter int, actor Actor, scenario Scenario) {
	for scenario.Scan() {
		actor.Act(scenario.Line())
	}
}
