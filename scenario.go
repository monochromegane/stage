package stage

type Scenario interface {
	Scan() bool
	Line() Line
}

type Line map[string]interface{}
