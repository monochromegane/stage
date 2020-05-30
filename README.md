# Stage [![Build Status](https://travis-ci.org/monochromegane/stage.svg?branch=master)](https://travis-ci.org/monochromegane/stage)

Simple and flexible simulation framework for Go.

This framework provides concurrent execution of simulations according to scenarios, output of results, monitor of progress and management of random number seeds.

All We need is implement the scenario and the actors who play it.

## Architecture

The algorithm of this framework is shown in the pseudo code below.

```go
for i := 0; i < iter; i++ {
    sem <- struct{}{}

    go func() {
        scenario := NewScenarioFn(rnd.Int63())
        actor := NewActorFn(rnd.Int63())

        for scenario.Scan() {
            action, _ := actor.Act(scenario.Line())
            w.Write(action.String())
        }
        <-sem
    }()
}
```

## Usage

1. Implement our scenario.
1. Implement an actor who plays the scenario.
1. Implement an action which represents performance of the actor.
1. Implement a callback function which is called at each iteration. (Optional)
1. And stage.New().Run()

```go
dir := "."                      // Directory for output
concurrency := runtime.NumCPU() // Number of concurrency for scenario
seed := 1                       // Seed for random
iter := 10                      // Number of iteration

s := stage.New(dir, concurrency, seed)
s.Run(iter, NewActorFn, NewScenarioFn, stage.NoOpeCallbackFn)
```

See also [examples](https://github.com/monochromegane/stage/blob/master/_examples).


### Scenario

Scenario interface represents our simulation scenario.
This framework runs the scenario number of iterations times in parallel.

Our scenario has Scan() method.
This method reads the scenario one line.
So, we can implements it as a counter or file reader.

```go
type Scenario struct {
    t     int
    limit int
    rnd   *rand.Rand
}

func (s *Scenario) Scan() bool {
    s.t++
    return s.t < s.limit
}
```

And the scenario has Line() method too.
This method returns the current line for the scenario as `stage.Line`.
Actors will perform according to it by each line.
We can use `stage.Line` flexibly because the struct is type of `map[string]interface{}`.

```go
func (s *Scenario) Line() stage.Line {
    return stage.Line{"x": s.rnd.NormFloat64()}
}
```

Finally, we define a method to generate the scenario for each iteration.

```go
func NewScenarioFn(seed int64) stage.Scenario {
    return &Scenario{
        t:     -1,
        limit: 3000,
        rnd:   rand.New(rand.NewSource(seed)),
    }
}
```

### Actor and Action

Actor interface represents our simulation actor.
This framework runs the actor number of lines times from the scenario.

Our actor has Act() method.
This method is called with a line of scenario one by one and returns a result of simulation as `stage.Action`.

```go
type Actor struct {
    n   int
    sum int
}

func (s *Actor) Act(line stage.Line) (stage.Action, error) {
    s.n++
    s.sum += line["x"].(float64)
    return Action{avg: float64(s.sum/s.n)}
}
```

We can implement `stage.Action` like a `Stringer`.

```go
type Action struct {
    avg float64
}

func (a Action) String() string {
    return fmt.Sprintf("%f\n", a.avg)
}
```

Finally, we define a method to generate the actor for each iteration.

```go
func NewActorFn(seed int64) stage.Actor {
    return &Actor{}
}
```

### Callback function (Optional)

If our scenario is going to run long, we can monitor it's progress by using `stage.CallbackFn`.
The callback function is called when each iteration finished.
So, we can use our favorite progress bar library.

```go
callbackFn := func(i int) { bar.Increment() } }
```

The framework has also an empty operation callback function named `stage.NoOpeCallbackFn`.
We usually use the function if we don't need monitor of progress.

## Installation

```sh
$ go get github.com/monochromegane/stage
```

## License

[MIT](https://github.com/monochromegane/stage/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
