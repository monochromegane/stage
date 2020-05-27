# Stage [![Build Status](https://travis-ci.org/monochromegane/stage.svg?branch=master)](https://travis-ci.org/monochromegane/stage)

Simple and flexible simulation framework for Go.

This framework provides concurrent execution of simulations according to scenarios, output of results, and management of random number seeds.

All We need is implement the scenario and the actors who play it.

## Usage

```go
dir := ""                       // Directory for output
concurrency := runtime.NumCPU() // Number of concurrency for scenario
seed := 1                       // Seed for random
iter := 10                      // Number of iteration

s := stage.New(dir, concurrency, seed)
s.Run(iter, NewYourActorFn, NewYourScenarioFn, stage.NoOpeCallbackFn)
```

The stage needs a scenario and an actor.

```go
type YourScenario struct{} // Implemantaion of stage.Scenario which has Scan and Line methods.

func NewYourScenarioFn(seed int64) Scenario {
    return &YourScenario{}
}

type YourActor struct{} // Implemantaion of stage.Actor which has Act method.

func NewYourActorFn(seed int64) Actor {
    return &YourActor{}
}
```

The actor return an action as Stringer.

```go
type YourAction struct{} // Implemantaion of Stringer which has String method.
```

If the show is too long, we can receive progress using callback function.

```go
func YourCallbackFn(i int) {} // Implemantaion of stage.CallbackFn method.
```

## Installation

```sh
$ go get github.com/monochromegane/stage
```

## License

[MIT](https://github.com/monochromegane/stage/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
