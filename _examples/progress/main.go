package main

import (
	"flag"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/monochromegane/stage"
)

var (
	i int
	o string
	s float64
)

func init() {
	flag.IntVar(&i, "i", 200, "Number of iteration")
	flag.StringVar(&o, "o", "example_log", "Directory name for result output")
}

func main() {
	flag.Parse()

	concurrency := -1 // Auto
	seed := int64(1)

	bar := pb.StartNew(i)
	callbackFn := func(i int) { bar.Increment() }
	err := stage.New(o, concurrency, seed).Run(i, NewActorFn, NewScenarioFn, callbackFn)
	if err != nil {
		panic(err)
	}
	bar.Finish()
}

// Scenario
type Scenario struct {
	t     int
	limit int
	s     time.Duration
}

func (s *Scenario) Scan() bool {
	s.t++
	return s.t < s.limit
}

func (s *Scenario) Line() stage.Line {
	return stage.Line{"s": s.s}
}

func NewScenarioFn(seed int64) stage.Scenario {
	return &Scenario{
		t:     -1,
		limit: 1,
		s:     time.Millisecond * 300,
	}
}

// Actor
type Actor struct{}

func (a *Actor) Act(line stage.Line) (stage.Action, error) {
	s := line["s"].(time.Duration)
	time.Sleep(s)
	return Action{}, nil
}

func NewActorFn(seed int64) stage.Actor {
	return &Actor{}
}

// Action
type Action struct{}

func (a Action) String() string {
	return ""
}
