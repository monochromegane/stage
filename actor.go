package stage

type NewActorFn func(seed int64) Actor

type Actor interface {
	Act(Line) (Action, error)
}

type Action interface {
	String() string
}
