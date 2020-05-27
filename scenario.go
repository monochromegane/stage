package stage

type NewScenarioFn func(seed int64) Scenario

type Scenario interface {
	Scan() bool
	Line() Line
}

type Line map[string]interface{}
