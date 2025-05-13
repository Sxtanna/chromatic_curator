package common

type Group []*actor

type actor = struct {
	Execute   func() error
	Interrupt func(error)
}

func (g Group) Act(execute func() error, interrupt func(error)) Group {
	return append(g, &actor{execute, interrupt})
}

func (g Group) Await(abort <-chan struct{}) {
	select {
	case <-abort:
	}
}
