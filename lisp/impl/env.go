package lisp

type Frame []Object
type Env struct {
	frame Frame
	next  *Env
}

type Location struct {
	level, offset int
}

func (env *Env) Push(frame Frame) *Env {
	return &Env{frame, env}
}

func (env *Env) Pop() *Env {
	if env == nil {
		panic("env underflow")
	}
	return env.next
}

func (env *Env) Locate(loc *Location) Object {
	for i := 0; i < loc.level; i++ {
		if env == nil {
			panic("illegal access to lexical environment")
		}
		env = env.next
	}
	frame := env.frame
	if loc.offset >= len(frame) {
		panic("illegal access to lexical environment")
	}
	return frame[loc.offset]
}
