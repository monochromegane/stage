package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/monochromegane/adwin"
	"github.com/monochromegane/stage"
)

var (
	i int
	o string
	s string
)

func init() {
	flag.IntVar(&i, "i", 1, "Number of iteration")
	flag.StringVar(&o, "o", "example_log", "Directory name for result output")
	flag.StringVar(&s, "s", "abrupt", "Type of scenario [abrupt|gradual]")
}

func main() {
	flag.Parse()

	concurrency := -1 // Auto
	seed := int64(1)
	newScenarioFn := NewScenarioAbruptChangeFn
	if s == "gradual" {
		newScenarioFn = NewScenarioGradualChangeFn
	}

	err := stage.New(o, concurrency, seed).Run(i, NewActorAdwinFn, newScenarioFn, stage.NoOpeCallbackFn)
	if err != nil {
		panic(err)
	}
}

// Scenario (Abrupt changes)
type ScenarioAbruptChange struct {
	t     int
	limit int
	rnd   *rand.Rand
}

func (s *ScenarioAbruptChange) Scan() bool {
	s.t++
	return s.t < s.limit
}

func (s *ScenarioAbruptChange) Line() stage.Line {
	std := 0.05
	mu := 0.8
	if s.t > 1000 {
		mu = 0.4
	}
	x := s.rnd.NormFloat64()*std + mu
	return stage.Line{
		"mu": mu,
		"x":  x,
	}
}

func NewScenarioAbruptChangeFn(seed int64) stage.Scenario {
	return &ScenarioAbruptChange{
		t:     -1,
		limit: 3000,
		rnd:   rand.New(rand.NewSource(seed)),
	}
}

// Scenario (Gradual changes)
type ScenarioGradualChange struct {
	t     int
	limit int
	rnd   *rand.Rand
}

func (s *ScenarioGradualChange) Scan() bool {
	s.t++
	return s.t < s.limit
}

func (s *ScenarioGradualChange) Line() stage.Line {
	std := 0.05
	mu := 0.0
	if s.t < 1000 {
		mu = 0.8
	} else if s.t < 2000 {
		mu = -0.0006*float64(s.t-1000) + 0.8
	} else {
		mu = 0.2
	}
	x := s.rnd.NormFloat64()*std + mu
	return stage.Line{
		"mu": mu,
		"x":  x,
	}
}

func NewScenarioGradualChangeFn(seed int64) stage.Scenario {
	return &ScenarioGradualChange{
		t:     -1,
		limit: 4000,
		rnd:   rand.New(rand.NewSource(seed)),
	}
}

// Actor
type ActorAdwin struct {
	adwin *adwin.Adwin
}

func (a *ActorAdwin) Act(line stage.Line) (stage.Action, error) {
	mu := line["mu"].(float64)
	x := line["x"].(float64)
	a.adwin.Add(x)

	return ActionAdwin{
		x:    x,
		mu:   mu,
		sum:  a.adwin.Sum(),
		size: a.adwin.Size(),
	}, nil
}

func NewActorAdwinFn(seed int64) stage.Actor {
	return &ActorAdwin{
		adwin: adwin.NewAdwin(0.01),
	}
}

// Action
type ActionAdwin struct {
	x    float64
	mu   float64
	sum  float64
	size int
}

func (a ActionAdwin) String() string {
	return fmt.Sprintf("%f,%f,%f,%d\n", a.x, a.mu, a.sum, a.size)
}
