package stage

type NewActorFn func() Actor

type Actor interface {
	Act(Line) (Action, error)
}

type Action interface {
	String() string
}
