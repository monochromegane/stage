# Stage [![Build Status](https://travis-ci.org/monochromegane/stage.svg?branch=master)](https://travis-ci.org/monochromegane/stage)

Simple and flexible simulation framework for Go.

This framework provides concurrent execution of simulations according to scenarios, output of results, monitor of progress and management of random number seeds.

All We need is implement the scenario and the actors who play it.

## Usage

1. Implement your scenario.
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
s.Run(iter, NewYourActorFn, NewYourScenarioFn, stage.NoOpeCallbackFn)
```

See also [examples](https://github.com/monochromegane/stage/blob/master/_examples).

### Scenario

Scenario interface represents our simulation scenario.
The stage package runs the scenario number of iterations times in parallel.

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
We can use `stage.Line` flexible because the struct is type of `map[string]interface{}`.

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

### Actor
### Action
### Callback function (Optional)

## Installation

```sh
$ go get github.com/monochromegane/stage
```

## License

[MIT](https://github.com/monochromegane/stage/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
