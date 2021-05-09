package lisp

type Frame []Object
type Env []Frame

type Location struct {
	level, offset int
}

func (env Env) Push(frame Frame) Env {
	return append(env, frame)
}

func (env Env) Pop() Env {
	if len(env) == 0 {
		panic("env underflow")
	}
	return env[:len(env)-1]
}

func (env Env) Locate(loc *Location) Object {
	if loc.level >= len(env) {
		panic("illegal access to lexical environment")
	}
	frame := env[loc.level]
	if loc.offset >= len(frame) {
		panic("illegal access to lexical environment")
	}
	return frame[loc.offset]
}
