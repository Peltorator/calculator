package eval

import "strconv"

type ExpressionEvaluator interface {
	Evaluate(expr string) (string, error)
}

type IncrementalFakeEvaluator struct {
	Current int
}

func (ev *IncrementalFakeEvaluator) Evaluate(expr string) (string, error) {
	result := strconv.Itoa(ev.Current)
	ev.Current++
	return result, nil
}