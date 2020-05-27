package stage

type NewScenarioFn func() Scenario

type Scenario interface {
	Scan() bool
	Line() Line
}

type Line map[string]interface{}
