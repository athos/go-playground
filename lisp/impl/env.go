package lisp

import "errors"

type Frame []Object
type Env struct {
	frame Frame
	next *Env
}

type Location struct {
	level, offset int
}

func NewEnv() *Env {
	return nil
}

func Push(env *Env, frame Frame) *Env {
	return &Env{frame, env}
}

func (env *Env) Pop() (*Env, error) {
	if env == nil {
		return nil, errors.New("Env underflow")
	}
	return env.next, nil
}

func (env *Env) Locate(loc *Location) (Object, error) {
	var frame Frame
	found := false
	for i := 0; i < loc.level; i++ {
		if env == nil {
			return nil, errors.New("illegal access to lexical environment")
		}
		frame = env.frame
		env = env.next
	}
	if !found {
		return nil, errors.New("illegal access to lexical environment")
	}
	if loc.offset >= len(frame) {
		return nil, errors.New("illegal access to lexical environment")
	}
	return frame[loc.offset], nil
}
