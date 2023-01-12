package query

import (
	"errors"
)

// 栈
type Stack struct {
	stack []string
}

func NewStack() *Stack {
	s := make([]string, 0)
	return &Stack{stack: s}
}

func (s *Stack) Pop() (string, error) {
	if s.isEmpty() {
		return "", errors.New("grammatical error")
	}
	value := s.stack[len(s.stack)-1]
	ns := make([]string, 0)
	for i := 0; i < len(s.stack)-1; i++ {
		ns = append(ns, s.stack[i])
	}
	s.stack = ns
	return value, nil
}

func (s *Stack) Push(value string) {
	ns := make([]string, 0)
	ns = append(ns, s.stack...)
	ns = append(ns, value)
	s.stack = ns
}

func (s *Stack) isEmpty() bool {
	if len(s.stack) == 0 {
		return true
	}
	return false
}
